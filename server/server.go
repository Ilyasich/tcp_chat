package main

import (
	"fmt"
	"log"
	"net"
	"tcp_chat/config"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	fmt.Println("Server has started!")
	if err != nil {
		log.Fatal(err)
	}
	go config.Broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go config.HandleConn(conn)
	}
}
