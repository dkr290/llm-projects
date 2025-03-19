package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

// Document represents a document to be stored in OpenSearch.
type Document struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}

// StoreDocumentBatch stores a batch of documents in OpenSearch.
func StoreDocumentBatch(
	ctx context.Context,
	client *qdrant.Client,
	indexName string,
	documents []Document,
) (time.Duration, error) {
	start := time.Now()
	points := make([]*qdrant.PointStruct, 0, len(documents))
	for _, doc := range documents {
		point := &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{
					Num: uint64(doc.ID),
				}, // Use document ID as the point ID
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: doc.Embedding, // Embedding (vector)
					},
				},
			},
			Payload: map[string]*qdrant.Value{
				"text": {Kind: &qdrant.Value_StringValue{StringValue: doc.Text}},
			},
		}
		points = append(points, point)
	}
	// Create an UpsertPoints request
	request := &qdrant.UpsertPoints{
		CollectionName: indexName,
		Points:         points,
	}

	// Perform the upsert operation
	_, err := client.Upsert(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("failed to perform bulk insert: %w", err)
	}
	return time.Since(start), nil // Return the time taken for this batch
}

// GenerateDocuments generates a list of documents for testing.
func GenerateDocuments(numDocs int) []Document {
	documents := make([]Document, numDocs)
	for i := range numDocs {

		embedding := make([]float32, 768)
		for j := range 768 {
			embedding[j] = rand.Float32()
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
	client *qdrant.Client,
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
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "qdrant",
		Port: 6334,
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

	CallBulkRead(ctx)
}

func CallBulkRead(ctx context.Context) {
	indexName := "test_collection"
	ids := make([]*qdrant.PointId, 1000)
	for i := 0; i < 1000; i++ {
		// Create a PointId with a random numeric ID
		ids[i] = &qdrant.PointId{
			PointIdOptions: &qdrant.PointId_Num{
				Num: uint64(rand.Intn(100000) + 1), // Random IDs between 1 and 1,000,000
			},
		}
	}

	conn, err := grpc.Dial("127.0.0.1:6334", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Qdrant: %v", err)
	}
	defer conn.Close()
	client1 := qdrant.NewPointsClient(conn)
	// Measure bulk read performance
	readTime, err := BulkReadDocuments(ctx, client1, indexName, ids)
	if err != nil {
		log.Fatalf("Error reading documents: %v", err)
	}

	log.Printf("Bulk read of 1000 documents completed in: %v", readTime)
}

func BulkReadDocuments(
	ctx context.Context,
	client qdrant.PointsClient,
	collectionName string,
	ids []*qdrant.PointId,
) (time.Duration, error) {
	start := time.Now() // Start time for the bulk read operation

	// Prepare the request to get points by IDs
	request := &qdrant.GetPoints{
		CollectionName: collectionName,
		Ids:            ids,
	}

	// Perform the get points request
	_, err := client.Get(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("failed to perform bulk read: %w", err)
	}

	return time.Since(start), nil // Return the time taken for the bulk read
}
