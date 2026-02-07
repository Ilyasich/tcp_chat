package server

import (
	"fmt"
	"log"
	"net"

	"git/config"

	"go.uber.org/zap"
)

func StartServer() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	lg := logger.Sugar()

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		lg.Fatal(err)
	}

	fmt.Println("Server has started!")
	lg.Info("Server listening on localhost:8000")

	// goroutine that will handle messages and connected clients
	go config.Broadcaster(lg)

	for {
		conn, err := listener.Accept()
		if err != nil {
			lg.Warnw("accept error", "error", err)
			continue
		}
		go config.HandleConn(lg, conn)
	}
}
