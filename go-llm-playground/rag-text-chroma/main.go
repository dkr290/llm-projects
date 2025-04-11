package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rag-text-chroma/pkg/chroma"

	"github.com/joho/godotenv"
)

func main() {
	apiKey := loadEnv()
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
		apiKey,
	)
	err = c.CreateEmbeddings()
	if err != nil {
		log.Println(err)
		return
	}
}

func loadEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	return apiKey
}
