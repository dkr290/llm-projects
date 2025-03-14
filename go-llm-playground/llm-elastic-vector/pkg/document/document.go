package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Document struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
}

type SearchResult struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
	Score     float64   `json:"score"`
}

func StoreDocument(ctx context.Context, client *opensearch.Client, indexName string) error {
	// Insert documents.

	d1 := Document{
		ID:        1,
		Text:      "this is text 1",
		Embedding: []float64{1.0, 2.0, 3.0, 4.0},
	}

	d2 := Document{
		ID:        2,
		Text:      "this is text 2",
		Embedding: []float64{1.5, 2.5, 3.5, 4.5},
	}
	documents := []Document{d1, d2}
	for _, doc := range documents {
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		req := opensearchapi.IndexRequest{
			Index:      indexName,
			DocumentID: fmt.Sprintf("%d", doc.ID), // Use document ID as the OpenSearch document ID
			Body:       strings.NewReader(string(docJSON)),
		}
		res, err := req.Do(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("error inserting document: %s", res.String())
		}

	}
	log.Println("Documents inserted successfully")
	return nil
}
