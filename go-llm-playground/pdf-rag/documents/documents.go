package documents

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func PdfLoader(path string, ctx context.Context) ([]schema.Document, error) {
	// Open the PDF file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening PDF: %v", err)
	}
	defer f.Close()
	FileInfo, _ := f.Stat()

	// Create a PDF loader
	loader := documentloaders.NewPDF(f, FileInfo.Size()) // No-op closer since defer handles it

	// Load and split the document
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 1200
	split.ChunkOverlap = 300
	chunks, err := loader.LoadAndSplit(ctx, split)
	if err != nil {
		return nil, fmt.Errorf("error loading and splitting PDF: %v", err)
	}

	fmt.Println("Number of chunks:", len(chunks))

	return chunks, nil
}
