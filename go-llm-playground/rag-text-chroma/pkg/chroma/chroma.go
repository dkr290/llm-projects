package chroma

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func TextToChunksDb(filePath string) ([]schema.Document, error) {
	// Construct the file path to the text file

	fmt.Println("Persistent directory does not exist. Initializing vector store...")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("The file %s does not exist. Please check the path.", filePath))
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %v", err)
	}
	defer file.Close()

	t := documentloaders.NewText(file)
	textSplitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1000),
		textsplitter.WithChunkOverlap(50),
	)
	doc, err := t.LoadAndSplit(context.Background(), textSplitter)
	if err != nil {
		return nil, fmt.Errorf("error the documensplitter %v", err)
	}
	fmt.Println("--- Document Chunks Information ---")
	fmt.Println("Number of document chunks:", len(doc))
	fmt.Println("Sample chunk:", doc[0].PageContent)

	return doc, nil
}
