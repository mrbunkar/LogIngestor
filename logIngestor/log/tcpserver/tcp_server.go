package tcpserver

import (
	"fmt"
	"log"
	"logIngestor/logIngestor/log/database"
	"net"
)

type TcpServer struct {
	ListenAddr string
	DbClient   database.ClientDB
	filter     chan string
	quitch     chan struct{}
}

func NewTCPServer(listenAddr string, mongoClient *database.MongoClient) *TcpServer {
	return &TcpServer{
		ListenAddr: listenAddr,
		DbClient:   mongoClient,
		filter:     make(chan string),
		quitch:     make(chan struct{}),
	}
}

func (ts *TcpServer) Close() error {
	fmt.Println("Closing the TCP Server...")
	return nil
}

func (ts *TcpServer) Start() error {
	listener, err := net.Listen("tcp", ts.ListenAddr)
	log.Println("Starting TCP Server on:", ts.ListenAddr)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go ts.handleTCPConnection(conn)

	}
}

func (ts *TcpServer) loop() {

	for {
		select {
		case message := <-ts.filter:
			ts.getLogs(message)
		}
	}
}

func (ts *TcpServer) handleTCPConnection(conn net.Conn) error {

	defer func() {
		fmt.Println("Dropping Client Connection")
		conn.Close()
	}()

	for {

	}
}

func (ts *TcpServer) addLogs(message string) {}

func (ts *TcpServer) getLogs(message string) {}
