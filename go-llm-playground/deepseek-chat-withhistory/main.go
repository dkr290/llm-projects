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
	loadEnv()
	ch := chat.New("You are a helpful assistant that explains complex topics step by step")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := ch.DsChat(ctx); err != nil {
		log.Println(err)
		return
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	os.Setenv("DEEPSEEK_API_KEY", apiKey)
}
