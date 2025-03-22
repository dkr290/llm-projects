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

func CreateIndex(index string, client *opensearch.Client) error {
	// Define the index settings and mappings
	indexSettings := map[string]any{
		"settings": map[string]any{
			"index": map[string]any{
				"knn":                true, // Enable k-NN search
				"number_of_shards":   2,    // Number of primary shards
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

	existsIndex := opensearchapi.IndicesExistsRequest{
		Index: []string{index},
	}
	existsRes, err := existsIndex.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error checking index existence: %s", err)
	}
	if existsRes.StatusCode == 200 {
		fmt.Println("The index exists, deleting the index")
		toDel := opensearchapi.IndicesDeleteRequest{
			Index: []string{index},
		}
		_, err := toDel.Do(context.Background(), client)
		if err != nil {
			return fmt.Errorf("error deleting tghe index %v", err)
		}
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
