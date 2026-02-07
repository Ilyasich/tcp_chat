package config

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
)

type client chan<- string

var (
	entering = make(chan client) // канал для регистрации новых клиентов
	leaving  = make(chan client) // канал для отслеживания клиентов, которые покидают чат
	messages = make(chan string) // канал для передачи сообщений между клиентами
)

func Broadcaster(lg *zap.SugaredLogger) {
	clients := make(map[client]bool) // мапа для отслеживания подключенных клиентов
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering: // получение нового клиента
			clients[cli] = true
		case cli := <-leaving: // клиент покидает чат
			delete(clients, cli)
			close(cli) // закрываем канал
		}
	}
}

// HandleConn обрабатывает подключение клиента к серверу
func HandleConn(lg *zap.SugaredLogger, conn net.Conn) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	conn.Write([]byte(fmt.Sprintf("Welcome to chat!\nDate: %s\n\nEnter your nickname: ", currentTime)))

	reader := bufio.NewReader(conn)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	fmt.Println("New user:", nickname)
	lg.Infof("New user: %s", nickname)

	conn.Write([]byte(fmt.Sprintf("Welcome, %s! Date: %s\n", nickname, currentTime)))

	ch := make(chan string)
	go ClientWriter(conn, ch)
	ch <- "You are: " + nickname
	messages <- nickname + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- nickname + ": " + input.Text()
	}

	leaving <- ch
	messages <- nickname + " has left chat"
	fmt.Println(nickname + " has left chat")
	lg.Infof("%s has left chat", nickname)
	conn.Close()
}

// ClientWriter отправляет все сообщения из канала клиенту
func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
