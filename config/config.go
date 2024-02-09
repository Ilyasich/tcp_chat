package config

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func Broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func HandleConn(conn net.Conn) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	conn.Write([]byte(fmt.Sprintf("Welcome to chat! \nData: %s\n \nEnter your nickname:", currentTime)))
	reader := bufio.NewReader(conn)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)
	fmt.Println("New user:", nickname)

	conn.Write([]byte(fmt.Sprintf("Welcome! Data: %s\n", currentTime)))

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
	conn.Close()
}

func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
