package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Document struct {
	ID        int       `bson:"id"`
	Text      string    `bson:"text"`
	Embedding []float64 `bson:"embedding"`
}

type SearchResult struct {
	ID        int       `bson:"id"`
	Text      string    `bson:"text"`
	Embedding []float64 `bson:"embedding"`
	Score     float64   `bson:"score"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Set the Server API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// Set client options
	opts := options.Client().
		ApplyURI("mongodb://172.22.0.2:27017/?directConnection=true").
		SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}
	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	err = CreateCollection(client)
	if err != nil {
		log.Fatalln(err)
	}
	err = CreateVectorIndex(client)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
