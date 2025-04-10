package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rag-text-chroma/pkg/chroma"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(currentDir, "documents", "lord_of_the_rings.txt")
	// Construct the path to the persistent directory

	d, err := chroma.TextToChunksDb(filePath)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	c := chroma.NewEmbeddings(
		d,
		"deepseek-chat",
		"deepseek-r1",
		"https://api.deepseek.com/v1",
		"http://localhost:8000",
		"sample_ns",
		token,
	)
}
