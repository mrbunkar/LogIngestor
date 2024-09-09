package database

import (
	"context"
	"fmt"
	"log"
	"logIngestor/logIngestor/log/logtype"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client DB can be usedas interface for others type of DB also, MySql, PostgreSql
type ClientDB interface {
	AddOne(*logtype.Log) error
	GetOne(any) (*logtype.Log, error)
	Close() error
}

type MongoConfig struct {
	User       string
	Password   string
	Collection string
	Database   string
	URI        string
}

type MongoClient struct {
	client *mongo.Client
	config *MongoConfig
	mu     sync.Mutex
}

func GetMongoClient() (*MongoClient, error) {
	mongoConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load MongoDB configuration: %w", err)
	}

	clientOpts := options.Client().ApplyURI(mongoConfig.URI).SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		client.Disconnect(context.TODO())
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}
	return &MongoClient{
		client: client,
		config: mongoConfig,
		mu:     sync.Mutex{},
	}, nil
}

func (mc *MongoClient) AddOne(lg *logtype.Log) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	collection := mc.client.Database(mc.config.Database).Collection(mc.config.Collection)
	_, err := collection.InsertOne(context.TODO(), lg)
	if err != nil {
		return fmt.Errorf("failed to insert log: %v", err)
	}
	// @TODO: return document Id also
	log.Println("Document added to Database")
	return nil
}

func (mc *MongoClient) AddMany(logs []*logtype.Log) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	collection := mc.client.Database(mc.config.Database).Collection(mc.config.Collection)

	var interfaceSlice []interface{}
	for _, log := range logs {
		interfaceSlice = append(interfaceSlice, log)
	}

	_, err := collection.InsertMany(context.TODO(), interfaceSlice)

	if err != nil {
		return fmt.Errorf("failed to insert log: %v", err)
	}
	// @TODO check for document Id also
	log.Println("All documents are insrted to Database")
	return nil
}

func (mc *MongoClient) GetOne(filter interface{}) (*logtype.Log, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	collection := mc.client.Database(mc.config.Database).Collection(mc.config.Collection)
	var log logtype.Log
	err := collection.FindOne(context.TODO(), filter).Decode(&log)
	if err != nil {
		return nil, fmt.Errorf("failed to find log: %v", err)
	}
	return &log, nil
}

func (mc *MongoClient) GetMany(filter interface{}) ([]*logtype.Log, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	collection := mc.client.Database(mc.config.Database).Collection(mc.config.Collection)

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find logs: %v", err)
	}
	defer cur.Close(context.Background())

	var results []*logtype.Log
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("failed to decode logs: %v", err)
	}

	return results, nil
}

func (mc *MongoClient) Close() error {
	return mc.client.Disconnect(context.TODO())
}

func LoadConfig() (*MongoConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return nil, err
	}

	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	collection := os.Getenv("COLLECTION")
	database := os.Getenv("DATABASE")

	if user == "" || password == "" {
		return nil, fmt.Errorf("MONGO_USER and MONGO_PASSWORD must be set")
	}

	return &MongoConfig{
		User:       user,
		Password:   password,
		Collection: collection,
		Database:   database,
		URI:        fmt.Sprintf("mongodb://%s:%s@localhost:27017", user, password),
	}, nil
}
