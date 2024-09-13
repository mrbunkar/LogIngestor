package logtype

import (
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Log struct {
	Level      string            `json:"level" bson:"level"`
	Message    string            `json:"message" bson:"message"`
	ResourceId string            `json:"resourceId" bson:"resourceId"`
	Timestamp  time.Time         `json:"timestamp" bson:"timestamp"` // Change to time.Time
	TraceId    string            `json:"traceId" bson:"traceId"`
	SpanId     string            `json:"spanId" bson:"spanId"`
	Commit     string            `json:"commit" bson:"commit"`
	Metadata   map[string]string `json:"metadata" bson:"metadata"`
}

func (l *Log) GetBsonEncoding() ([]byte, error) {
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
