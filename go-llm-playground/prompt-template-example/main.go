package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

func main() {
	token := loadEnv()
	// Initialize the OpenAI client with Deepseek model
	llm, err := openai.New(
		openai.WithModel("deepseek-chat"),
		openai.WithBaseURL("https://api.deepseek.com/v1"),
		openai.WithToken(token),
	)
	if err != nil {
		log.Fatalf("error connectint to deepseek %v", err)
	}
	messagess := []prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"You are a comedian who tells jokes about {{.topic}}.",
			[]string{"topic"},
		),
		prompts.NewHumanMessagePromptTemplate(
			"Tell me {{.joke_count}} jokes.",
			[]string{"joke_count"},
		),
	}
	promtTemplate := prompts.NewChatPromptTemplate(messagess)

	// Create the prompt values
	promptValues := map[string]any{
		"topic":      "lawyers",
		"joke_count": 3,
	}
	prmp, err := promtTemplate.Format(promptValues)
	if err != nil {
		log.Fatal(err)
	}

	// Invoke the LLM
	result, err := llm.Call(context.Background(), prmp)
	if err != nil {
		fmt.Println("Error calling LLM:", err)
		return
	}

	fmt.Println(result)
}

func loadEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	return apiKey
}
