package logtype

import (
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Log struct {
	Level      string            `json:"level" bson:"level"`
	Message    string            `json:"message" bson:"message"`
	ResourceId string            `json:"resource_id" bson:"resource_id"`
	Timestamp  string            `json:"timestamp" bson:"timestamp"`
	TraceId    string            `json:"trace_id" bson:"trace_id"`
	SpanId     string            `json:"span_id" bson:"span_id"`
	Commit     string            `json:"commit" bson:"commit"`
	Metadata   map[string]string `json:"metadata" bson:"metadata"`
}

func (l *Log) GetBson() ([]byte, error) {
	bsonData, err := bson.Marshal(l)
	if err != nil {
		log.Println("Failed to marshal log entry to BSON:", err)
		return nil, err
	}
	return bsonData, nil
}

func (l *Log) GetJsonEncoding() []byte {

	b, err := json.Marshal(l)
	if err != nil {
		log.Println(err)
	}
	return b
}

func NewLog(raw_data []byte) *Log {
	l := new(Log)
	err := json.Unmarshal(raw_data, l)
	if err != nil {
		log.Println(err)
	}
	return l
}
