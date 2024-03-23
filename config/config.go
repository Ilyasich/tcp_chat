package config

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type client chan<- string

var db *sql.DB

var (
	entering = make(chan client) //канал для регистрации новых клиентов
	leaving  = make(chan client) //канал для отслеживания клиентов, которые покидают чат
	messages = make(chan string) //канал для передачи сообщений между клиентами
	data     = make(chan string) //канал для бд ????
)

func Broadcaster(lg *zap.SugaredLogger) {
	clients := make(map[client]bool) //мапа для отслеживания подключенных клиентов
	//обработка событий
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering: //получение нового клиента
			clients[cli] = true
		//case cli := <-data:
		//SaveUserToDB(cli)//???????????
		case cli := <-leaving: //клиент покидает чат
			delete(clients, cli)
			close(cli) //закрываем канал
		}
	}
}


// обрабатывает подключение клиента к серверу
func HandleConn(lg zap.SugaredLogger, conn net.Conn) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	conn.Write([]byte(fmt.Sprintf("Welcome to chat! \nData: %s\n \nEnter your nickname:", currentTime))) //вывод на экран времени подключенного клиента
	reader := bufio.NewReader(conn)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)
	fmt.Println("New user:", nickname)
	lg.Info("New user:", nickname)

	conn.Write([]byte(fmt.Sprintf("Welcome! Data: %s\n", currentTime)))

	ch := make(chan string)   //канал для передачи информации о действии клиентов
	go ClientWriter(conn, ch) //горутина пишущая в канал никнеймы сообщения и кто присоеденился и покинул чат
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
	lg.Info(nickname + "has left chat")
	conn.Close()
}

// принимает соединение клиента и канал ch, и отправляет все полученные сообщения из канала клиенту
func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

//функция для сохранения в бд????????
func SaveUserToDB(nickname string) {
	var err error
	_, err = db.ExecContext(context.Background(), `INSERT INTO users (id, nickname) VALUES ($1, $2)`, nickname)

	if err != nil {
		log.Panic(err)
	}
}
