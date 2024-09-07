package database

import (
	"context"
	"fmt"
	"log"
	"logIngestor/logIngestor/log/logtype"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDb(t *testing.T) {

	mongoClient, err := GetMongoClient()

	if err != nil {
		log.Panic(err)
	}

	if err := mongoClient.client.Ping(context.TODO(), nil); err != nil {
		log.Panic(err)
	}

	defer mongoClient.Close()
	assert.Equal(t, mongoClient.config.User, "mrbunkar")
	assert.Equal(t, mongoClient.config.Password, "changeme")
}

func TestMongoReadWrite(t *testing.T) {
	mongoClient, err := GetMongoClient()

	if err != nil {
		log.Panic(err)
	}

	if err := mongoClient.client.Ping(context.TODO(), nil); err != nil {
		log.Panic(err)
	}

	defer mongoClient.Close()

	payload := `{
		"level": "error",
		"message": "Failed to connect to DB",
		"resource_id": "server-1234",
		"timestamp": "2023-09-15T08:00:00Z",
		"trace_id": "abc-xyz-123",
		"span_id": "span-456",
		"commit": "5e5342f",
		"metadata": {
			"parentResourceId": "server-0987"
		}}`

	logEntry := logtype.NewLog([]byte(payload))
	err = mongoClient.AddOne(logEntry)

	if err != nil {
		assert.Error(t, err, "Failed to upload the data")
		return
	}

	filter := bson.M{"resource_id": logEntry.ResourceId}

	fetchedLog, err := mongoClient.GetOne(filter)
	if err != nil {
		assert.Error(t, err, "Failed to get the data")
		return
	}

	assert.Equal(t, logEntry, fetchedLog)
	fmt.Println("Log Entry: ", logEntry)
	fmt.Println("Fetched log:", fetchedLog)
	assert.Equal(t, logEntry.Message, fetchedLog.Message)
	assert.Equal(t, logEntry.ResourceId, fetchedLog.ResourceId)
	assert.Equal(t, logEntry.Level, fetchedLog.Level)
	assert.Equal(t, logEntry.Timestamp, fetchedLog.Timestamp)
	assert.Equal(t, logEntry.TraceId, fetchedLog.TraceId)
	assert.Equal(t, logEntry.SpanId, fetchedLog.SpanId)
	assert.Equal(t, logEntry.Commit, fetchedLog.Commit)
	assert.Equal(t, logEntry.Metadata["parentResourceId"], fetchedLog.Metadata["parentResourceId"])
}

func TestEnvConfig(t *testing.T) {
	config, err := LoadConfig()

	if err != nil {
		log.Panic(err)
	}

	assert.Equal(t, config.User, "mrbunkar")
	assert.Equal(t, config.Password, "changeme")
}
