package pdf

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

// PDF is a struct that represents a PDF document.
type PDF struct {
	Content  string
	Filepath string
}

func New(p string) *PDF {
	return &PDF{
		Filepath: p,
	}
}

func (p *PDF) ReadPdf() error {
	// Open the PDF file
	output, err := exec.Command("pdftotext", p.Filepath, "-").Output()
	if err != nil {
		return fmt.Errorf("failed to extract text: %v", err)
	}

	p.Content = string(output)
	return nil
}

func (p *PDF) SplitText() ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1200),
		textsplitter.WithChunkOverlap(200),
	)
	// Use the splitter to split the text into chunks
	chunks, err := splitter.SplitText(p.Content)
	if err != nil {
		return nil, fmt.Errorf("error splitting the text %v", err)
	}

	return chunks, nil
}

func (p *PDF) AddMetadata(chunks []string, docTitle string) ([]schema.Document, error) {
	type metadata struct {
		title  string
		author string
		date   string
	}

	var documents []schema.Document

	// Loop through each chunk
	for _, c := range chunks {
		// Create a metadata struct for the current chunk
		md := metadata{
			title:  docTitle,
			author: "US Business Bureau",
			date:   time.Now().Format("2006-01-02 15:04:05"), // Format the time
		}

		// Create a schema.Document for the current chunk
		doc := schema.Document{
			PageContent: c,
			Metadata: map[string]interface{}{
				"title":  md.title,
				"author": md.author,
				"date":   md.date,
			},
		} // Append the document to the slice
		documents = append(documents, doc)
	}

	return documents, nil
}

func (p *PDF) GenEmbeddings(
	chunks []schema.Document,
	modelName, ollamaUrl string,
	qdrantUrl, collectionName string,
) error {
	ollamaLLM, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaUrl),
	)
	if err != nil {
		return fmt.Errorf("cannot load ollama model %v", err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		return fmt.Errorf("new embedder error %v", err)
	}

	// Create a new Qdrant vector store.
	url, err := url.Parse(qdrantUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Create a Qdrant client
	client, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		return fmt.Errorf("failed to create Qdrant client: %v", err)
	}
	// Add documents to Qdrant

	_, err = client.AddDocuments(context.Background(), chunks)
	if err != nil {
		return fmt.Errorf("failed to add documents to Qdrant: %v", err)
	}
	return nil
}
