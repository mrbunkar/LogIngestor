package httpserver

import (
	"fmt"
	"log"
	"logIngestor/logIngestor/log/database"
	"logIngestor/logIngestor/log/logtype"
	"logIngestor/logIngestor/log/validate"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	StatusOK                  = http.StatusOK
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed
	StatusBadRequest          = http.StatusBadRequest
	StatusFailedDependency    = http.StatusFailedDependency
	StatusInternalServerError = http.StatusInternalServerError
)

type HTTPServer struct {
	ListenAddr string
	DbClient   database.ClientDB
	Decoder    *validate.Decoder
}

func NewHTTPServer(listenAddr string, mongoClient *database.MongoClient) *HTTPServer {

	return &HTTPServer{
		ListenAddr: listenAddr,
		DbClient:   mongoClient,
		Decoder:    validate.NewDecoder(),
	}
}

func (s *HTTPServer) Close() error {
	return s.DbClient.Close()
}

func (s *HTTPServer) initialiseRoutes() {
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/Add", s.AddLogHandler)
	http.HandleFunc("/Get", s.GetLogHandler)
}

func (s *HTTPServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	n, err := fmt.Fprintln(w, "Welcome to LogIngestor", StatusOK)
	if err != nil {
		http.Error(w, "Unkown Error"+err.Error(), StatusInternalServerError)
	}
	log.Printf("[%d] Bytes sent over the Netowrk\n", n)
}

func (s *HTTPServer) AddLogHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to add logs
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
		return
	}
	logEntry := &logtype.Log{}

	if err := s.Decoder.Decode(r.Body, logEntry); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), StatusBadRequest)
		return
	}

	if err := s.DbClient.AddOne(logEntry); err != nil {
		http.Error(w, "Failed to retrieve logs from the database", StatusFailedDependency)
		return
	}

	w.WriteHeader(StatusOK)
	fmt.Fprintf(w, "Document added to DB")
}

func (s *HTTPServer) GetLogHandler(w http.ResponseWriter, r *http.Request) {

	// Implement logic to retrieve logs
	// @TODO add a validator to validate the log filters

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		http.Error(w, "Filter parameter is required", StatusBadRequest)
		return
	}

	// Parse the filter into a BSON document
	// @TODO validator{
	filters := strings.Split(filter, "&")
	bsonFilter := bson.M{}
	for _, f := range filters {
		parts := strings.Split(f, ":")
		if len(parts) == 2 {
			bsonFilter[parts[0]] = parts[1]
		}
	}

	//}
	log, err := s.DbClient.GetOne(bsonFilter)
	if err != nil {
		http.Error(w, "Failed to read from MongoDB. Error: "+err.Error(), http.StatusFailedDependency)
		return
	}

	w.WriteHeader(StatusOK)
	fmt.Fprint(w, string(log.GetJsonEncoding()))
}

func (s *HTTPServer) Start() error {
	log.Printf("Starting HTTP server on %s\n", s.ListenAddr)
	s.initialiseRoutes()
	if err := http.ListenAndServe(s.ListenAddr, nil); err != nil {
		s.Close()
		log.Fatalf("Error starting server: %v", err)
	}
	return nil
}
