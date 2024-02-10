package server

import (
	"fmt"
	"log"
	"net"

	"tcp_chat/config"
)

func StartServer() {
	listener, err := net.Listen("tcp", "localhost:8000") //создаем TCP-сервер, прослушивающий подключения на локальном адресе и порту 8080

	fmt.Println("Server has started!") //сообщаем о старте сервера
	if err != nil {
		log.Fatal(err)
	}
	go config.Broadcaster() //горутина которая будет обрабатывать сообщения и управлять подключенными клиентами.
	for {
		conn, err := listener.Accept() //слушаем порт
		if err != nil {
			log.Print(err)
			continue
		}
		go config.HandleConn(conn)
	}
}
