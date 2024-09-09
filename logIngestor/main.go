package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

	time.Sleep(1 * time.Second)
	populate()
}

func populate() {
	levels := []string{"INFO", "ERROR", "WARN", "DEBUG"}
	messages := []string{"This is a message 1", "This is a message 2", "This is a message 3", "This is a debug message"}
	for j := 0; j < 10; j++ {
		logReq := map[string]interface{}{
			"Level":      levels[rand.Intn(4)],
			"Message":    messages[rand.Intn(4)],
			"ResourceId": generateRandomStringNumber(4),
			"Timestamp":  time.Now(),
			"TraceId":    "trace-" + generateRandomStringNumber(4),
			"SpanId":     "span-" + generateRandomStringNumber(4),
			"Commit":     generateRandomStringNumber(6),
			"Metadata": map[string]string{
				"ParentResourceId": generateRandomStringNumber(4),
			},
		}
		jsonLogReq, err := json.Marshal(logReq)
		if err != nil {
			panic("error marshalling log request")
		}

		res, err := http.Post("http://localhost:3000/add", "application/json", bytes.NewBuffer(jsonLogReq))
		if err != nil {
			panic("error posting log request")
		}
		log.Println("Status code: ", res.StatusCode)

		timeToSleep := time.Duration(rand.Intn(200)) * time.Millisecond

		time.Sleep(timeToSleep)
	}
}

func generateRandomStringNumber(n int) string {
	const letters = "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
