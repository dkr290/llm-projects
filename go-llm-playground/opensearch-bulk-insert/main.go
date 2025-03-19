package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Document represents a document to be stored in OpenSearch.
type Document struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
}

// StoreDocumentBatch stores a batch of documents in OpenSearch.
func StoreDocumentBatch(
	ctx context.Context,
	client *opensearch.Client,
	indexName string,
	documents []Document,
) (time.Duration, error) {
	start := time.Now()
	var bulkRequestBody strings.Builder

	for _, doc := range documents {
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal document: %w", err)
		}

		bulkRequestBody.WriteString(
			fmt.Sprintf(
				`{ "index" : { "_index" : "%s", "_id" : "%d" } }%s`,
				indexName,
				doc.ID,
				"\n",
			),
		)
		bulkRequestBody.WriteString(string(docJSON) + "\n")
	}

	req := opensearchapi.BulkRequest{
		Index: indexName,
		Body:  strings.NewReader(bulkRequestBody.String()),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return 0, fmt.Errorf("failed to perform bulk insert: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error in bulk insert: %s", res.String())
	}

	return time.Since(start), nil // Return the time taken for this batch
}

// GenerateDocuments generates a list of documents for testing.
func GenerateDocuments(numDocs int) []Document {
	documents := make([]Document, numDocs)
	for i := range numDocs {

		embedding := make([]float64, 768)
		for j := range 768 {
			embedding[j] = rand.Float64()
		}

		documents[i] = Document{
			ID:        i + 1,
			Text:      fmt.Sprintf("this is text %d", i+1),
			Embedding: embedding,
		}
	}
	return documents
}

// StoreDocumentsConcurrently stores documents concurrently using goroutines.
func StoreDocumentsConcurrently(
	ctx context.Context,
	client *opensearch.Client,
	indexName string,
	numDocs, batchSize, numWorkers int,
) error {
	documents := GenerateDocuments(numDocs)
	chunkSize := (numDocs + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	var totalBulkTime time.Duration
	var mu sync.Mutex // Mutex to protect shared variables
	totalStart := time.Now()

	for i := range numWorkers {
		go func(workerID int) {
			defer wg.Done()

			start := workerID * chunkSize
			end := start + chunkSize
			if end > numDocs {
				end = numDocs
			}

			for j := start; j < end; j += batchSize {
				batchEnd := j + batchSize
				if batchEnd > end {
					batchEnd = end
				}

				batch := documents[j:batchEnd]
				batchTime, err := StoreDocumentBatch(ctx, client, indexName, batch)
				if err != nil {
					log.Printf("Worker %d failed to insert batch: %v", workerID, err)
					return
				}
				mu.Lock()
				totalBulkTime += batchTime // Accumulate bulk operation time
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(totalStart) // Total time for the entire operation

	log.Printf("All documents inserted successfully")
	log.Printf("Total time taken: %v", totalTime)
	log.Printf("Total time spent in bulk operations: %v", totalBulkTime)
	log.Printf(
		"Average time per bulk operation: %v",
		totalBulkTime/time.Duration(numDocs/batchSize),
	)
	return nil
}

func main() {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"https://opensearch-api.k8s-dev"},
		Username:  "admin",
		Password:  "",
	})
	if err != nil {
		log.Fatalf("Error creating OpenSearch client: %v", err)
	}

	ctx := context.Background()
	indexName := "test_index"

	numDocs := 100000
	batchSize := 100
	numWorkers := 5

	if err := StoreDocumentsConcurrently(ctx, client, indexName, numDocs, batchSize, numWorkers); err != nil {
		log.Fatalf("Error storing documents: %v", err)
	}

	CallBulkRead(ctx, client, indexName)
}

func CallBulkRead(ctx context.Context, client *opensearch.Client, indexName string) {
	ids := make([]string, 1000)
	for i := range 1000 {
		ids[i] = fmt.Sprintf("%d", rand.Intn(100000)+1) // Random IDs between 1 and 1,000,000
	}

	// Measure bulk read performance
	readTime, err := BulkReadDocuments(ctx, client, indexName, ids)
	if err != nil {
		log.Fatalf("Error reading documents: %v", err)
	}

	log.Printf("Bulk read of 1000 documents completed in: %v", readTime)
}

func BulkReadDocuments(
	ctx context.Context,
	client *opensearch.Client,
	indexName string,
	ids []string,
) (time.Duration, error) {
	start := time.Now() // Start time for the bulk read operation

	// Prepare the request body for multi-get
	var requestBody strings.Builder
	requestBody.WriteString(`{ "docs": [`)
	for i, id := range ids {
		if i > 0 {
			requestBody.WriteString(",")
		}
		requestBody.WriteString(fmt.Sprintf(`{ "_index": "%s", "_id": "%s" }`, indexName, id))
	}
	requestBody.WriteString(`] }`)

	// Perform the multi-get request
	req := opensearchapi.MgetRequest{
		Body: strings.NewReader(requestBody.String()),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return 0, fmt.Errorf("failed to perform bulk read: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error in bulk read: %s", res.String())
	}

	return time.Since(start), nil // Return the time taken for the bulk read
}
