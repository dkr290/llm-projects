package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func CreateIndex(index string, address string) error {
	// Create a client instance
	cfg := opensearch.Config{
		Addresses: []string{
			address, // Replace with your OpenSearch URL
		},
	}
	client, err := opensearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating the client: %s", err)
	}

	// Define the index settings and mappings
	indexSettings := map[string]any{
		"settings": map[string]any{
			"index": map[string]any{
				"knn":                true, // Enable k-NN search
				"number_of_shards":   3,    // Number of primary shards
				"number_of_replicas": 1,    // Number of replicas
			},
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"id": map[string]any{
					"type": "keyword",
				},
				"content": map[string]any{
					"type": "text",
				},
				"embedding": map[string]any{
					"type":       "knn_vector", // Use knn_vector for OpenSearch
					"dimension":  768,          // Number of dimensions
					"index":      true,
					"similarity": "cosine", // Similarity metric
				},
			},
		},
	}

	// Convert the settings to JSON
	settingsJSON, err := json.Marshal(indexSettings)
	if err != nil {
		log.Fatalf("Error marshaling settings to JSON: %s", err)
	}

	// Create the index
	req := opensearchapi.IndicesCreateRequest{
		Index: index, // Replace with your desired index name
		Body:  strings.NewReader(string(settingsJSON)),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error creating the index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response: %s", res.String())
	}

	fmt.Println("Index created successfully!")
	return nil
}
