package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"pdf-rag-qdrant/llmquery"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

func main() {
	// ctx := context.Background()
	// chunks, err := documents.PdfLoader("./data/how-to-code-in-go.pdf", ctx)
	// if err != nil {
	// 	slog.Error("erro chunking the file", "error", err)
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	// err = vector_embed(chunks, ctx)
	// if err != nil {
	// 	slog.Error("vector_embed error", "error", err)
	// }
	if err := llmquery.QuestionResponse(
		ctx,
	); err != nil {
		slog.Error("Error retreiving data: ", "error", err)
	}
}

func vector_embed(chunks []schema.Document, ctx context.Context) error {
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

	// Create a new Qdrant vector store.
	url, err := url.Parse("http://qdrant.172.22.0.3.nip.io/")
	if err != nil {
		log.Fatal(err)
	}
	store, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithCollectionName("vector_store"),
		qdrant.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		return err
	}

	_, err = store.AddDocuments(ctx, chunks)
	if err != nil {
		return err
	}
	return nil
}
