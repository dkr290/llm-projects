package main

import (
	"context"
	"deepseek-chat-withhistory/pkg/chat"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	apiKey := loadEnv()
	ch := chat.New("You are a helpful assistant that explains complex topics step by step")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := ch.DsChat(ctx, apiKey); err != nil {
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
