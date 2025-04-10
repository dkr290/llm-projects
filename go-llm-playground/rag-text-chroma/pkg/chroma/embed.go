package chroma

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type EmbedDocuments interface {
	ChromaEmbeddings() error
}

type ChromaEmbedding struct {
	doc             []schema.Document
	chatModel       string
	embedModel      string
	chatModelUrl    string
	token           string
	chromaUrl       string
	chromaNamespace string
}

func NewEmbeddings(
	doc []schema.Document,
	chatm, embedm, chaturl, chromaurl, chromans, token string,
) *ChromaEmbedding {
	return &ChromaEmbedding{
		doc:             doc,
		chatModel:       chatm,
		embedModel:      embedm,
		chatModelUrl:    chaturl,
		token:           token,
		chromaUrl:       chromaurl,
		chromaNamespace: chromans,
	}
}

func (c *ChromaEmbedding) CreateEmbeddings() error {
	fmt.Println("--- Creating embeddings ---")

	opts := []openai.Option{
		openai.WithModel(c.chatModel),
		openai.WithEmbeddingModel(c.embedModel),
		openai.WithBaseURL(c.chatModelUrl),
		openai.WithToken(c.token),
	}

	llm, err := openai.New(opts...)
	if err != nil {
		return fmt.Errorf("error connectint to deepseek %v", err)
	}

	embedings, err := llm.CreateEmbedding(context.Background(), c.doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- Creating vector store ---")
	_, err = chroma.New(
		chroma.WithChromaURL(c.chatModel), // or your Chroma server URL
		chroma.WithEmbedder(embedings),
		chroma.WithNameSpace(c.chromaNamespace), // optional namespace
	)
	if err != nil {
		panic(fmt.Errorf("failed to create vector store: %w", err))
	}
	return nil
}
