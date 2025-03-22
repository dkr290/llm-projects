package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	ops "github.com/opensearch-project/opensearch-go"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
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
			Metadata: map[string]any{
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
	opensearchClient *ops.Client, IndexName string,
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

	err = manualEmbed(chunks, ollamaEmbeder, opensearchClient, IndexName)
	if err != nil {
		return err
	}
	return nil
}

func manualEmbed(
	chunks []schema.Document,
	ollamaEmbeder *embeddings.EmbedderImpl,
	opensearchClient *ops.Client, IndexName string,
) error {
	// Alternative: Manually index documents using the OpenSearch client
	for i, chunk := range chunks {
		// Generate embedding for the chunk
		embedding, err := ollamaEmbeder.EmbedQuery(context.Background(), chunk.PageContent)
		if err != nil {
			return fmt.Errorf("failed to generate embedding for chunk %d: %v", i, err)
		}

		// Prepare the document for OpenSearch
		doc := map[string]any{
			"content":   chunk.PageContent,
			"embedding": embedding,
			// Add metadata or other fields as needed
		}
		// Step 1: Marshal the map to JSON
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document %d: %v", i, err)
		}

		// Step 2: Create an io.Reader from the JSON bytes
		docReader := bytes.NewReader(docJSON)

		// Index the document in OpenSearch
		_, err = opensearchClient.Index(
			IndexName,
			docReader,
			opensearchClient.Index.WithDocumentID(fmt.Sprintf("doc-%d", i)), // Optional: unique ID
			opensearchClient.Index.WithRefresh(
				"true",
			), // Ensure immediate visibility
		)
		if err != nil {
			return fmt.Errorf("failed to index document %d: %v", i, err)
		}
	}
	return nil
}
