package main

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const indexName = "vector_index"

// VectorIndexSettings represents setting to create a vector index.
type VectorIndexSettings struct {
	NumDimensions int
	Path          string
	Similarity    string
}

func CreateVectorIndex(client *mongo.Client) error {
	settings := VectorIndexSettings{
		NumDimensions: 4,
		Path:          "embedding",
		Similarity:    "cosine",
	}
	collection := client.Database(dbName).Collection(collectionName)
	return createVectorIndex(collection, settings)
}

func createVectorIndex(collection *mongo.Collection, settings VectorIndexSettings) error {
	// Define the vector search index model
	indexModel := mongo.SearchIndexModel{
		Definition: bson.D{
			{Key: "name", Value: indexName},
			{Key: "mappings", Value: bson.D{
				{Key: "dynamic", Value: true},
				{Key: "fields", Value: bson.D{
					{Key: settings.Path, Value: bson.D{
						{Key: "type", Value: "knnVector"},
						{Key: "dimensions", Value: settings.NumDimensions},
						{Key: "similarity", Value: settings.Similarity},
					}},
				}},
			}},
		},
	}
	// Create the search index
	_, err := collection.SearchIndexes().CreateOne(context.Background(), indexModel)
	return err
}
