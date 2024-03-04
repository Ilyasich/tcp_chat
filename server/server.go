package server

import (
	"git/config"

	"go.uber.org/zap"

	"fmt"
	"log"
	"net"

)

func StartServer() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	lg := logger.Sugar()

	listener, err := net.Listen("tcp", "localhost:8000")//создаем TCP-сервер, прослушивающий подключения на локальном адресе и порту 8080
	fmt.Println("Server has started!")
	lg.Info("Server has started...")
	if err != nil {
		lg.Fatal(err)
	}
	go config.Broadcaster(lg)//горутина которая будет обрабатывать сообщения и управлять подключенными клиентами.
	for {
		conn, err := listener.Accept()
		if err != nil {
			lg.Panicln("Port is not evelible", err)
			continue
		}
		go config.HandleConn(*lg, conn)
	}
}
