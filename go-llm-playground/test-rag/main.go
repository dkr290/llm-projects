package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"test-rag/llmquery"
	"time"

	ops "github.com/opensearch-project/opensearch-go"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/opensearch"
)

func main() {
	// Initialize OpenSearch client
	client, err := ops.NewClient(ops.Config{
		Addresses: []string{"http://api.172.22.0.4.nip.io"}, // OpenSearch server address
		Username:  "admin",
		Password:  "admin",
	})
	if err != nil {
		slog.Error("Error creating OpenSearch client: %s", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	//
	err = vector_embed(ctx, client)
	if err != nil {
		slog.Error("vector_embed error", "error", err)
	}
	if err := llmquery.QuestionResponse(
		ctx,
		"nomic-embed-text",
		"http://172.22.0.4/ollama",
		client); err != nil {
		slog.Error("Error retreiving data: ", "error", err)
	}
}

func vector_embed(ctx context.Context, client *ops.Client) error {
	ollamaLLM, err := ollama.New(
		ollama.WithModel("nomic-embed-text"),
		ollama.WithServerURL("http://172.22.0.4/ollama"),
	)
	if err != nil {
		return fmt.Errorf("cannot load ollama model %v", err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		return fmt.Errorf("new embedder error %v", err)
	}

	store, err := opensearch.New(client,
		opensearch.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		return err
	}
	_, err = store.CreateIndex(ctx, "vector_index")
	if err != nil {
		return fmt.Errorf("error creating index %v", err)
	}
	a := []schema.Document{
		{
			PageContent: "A city in texas",
			Metadata: map[string]any{
				"area": 3251,
			},
		},
		{
			PageContent: "A country in Asia",
			Metadata: map[string]any{
				"area": 2342,
			},
		},
		{
			PageContent: "A country in South America",
			Metadata: map[string]any{
				"area": 432,
			},
		},
		{
			PageContent: "An island nation in the Pacific Ocean",
			Metadata: map[string]any{
				"area": 6531,
			},
		},
		{
			PageContent: "A mountainous country in Europe",
			Metadata: map[string]any{
				"area": 1211,
			},
		},
		{
			PageContent: "A lost city in the Amazon",
			Metadata: map[string]any{
				"area": 1223,
			},
		},
		{
			PageContent: "A city in England",
			Metadata: map[string]any{
				"area": 4324,
			},
		},
	}
	resp, err := store.AddDocuments(ctx, a)
	if err != nil {
		return err
	}
	log.Println(resp)
	// Search for similar documents using score threshold.
	docs, err := store.SimilaritySearch(
		ctx,
		"american places",
		10,
		vectorstores.WithScoreThreshold(0.80),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(docs)
	return nil
}
