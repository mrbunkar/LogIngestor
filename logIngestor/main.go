package main

import (
	"log"
	"os"

	lg "logIngestor/logIngestor/log"
	"logIngestor/logIngestor/log/database"
	"logIngestor/logIngestor/log/httpserver"
	"logIngestor/logIngestor/log/tcpserver"
)

func startServer(server lg.Server, errChan chan error) {
	if err := server.Start(); err != nil {
		errChan <- err
	}
}

func main() {
	mongoClient, err := database.GetMongoClient()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	httpServer := httpserver.NewHTTPServer(":3000", mongoClient)
	tcpServer := tcpserver.NewTCPServer(":3030", mongoClient)

	httpErrChan := make(chan error)
	tcpErrChan := make(chan error)

	go startServer(httpServer, httpErrChan)
	go startServer(tcpServer, tcpErrChan)

	select {
	case err := <-httpErrChan:
		log.Printf("HTTP server stopped with error: %v", err)
		os.Exit(1)
	case err := <-tcpErrChan:
		log.Printf("TCP server stopped with error: %v", err)
		os.Exit(1)
	}
}
