package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	dbName         = "example4"
	collectionName = "book"
)

func CreateCollection(client *mongo.Client) error {
	db := client.Database(dbName)
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return fmt.Errorf("createCollection: %w", err)
	}
	fmt.Println("Created Collection")
	return nil
}
