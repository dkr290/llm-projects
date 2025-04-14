package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the GEMINI_API_KEY environment variable")
		return
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Define user and system prompts
	systemPrompt := "You are a helpful and concise assistant."
	userPrompt := "What is the capital of Germany?"

	model := "gemini-2.0-flash"

	// Construct the content with system and user messages
	userContent := []*genai.Part{
		{Text: userPrompt},
	}

	systemContent := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		[]*genai.Content{{Parts: userContent}},
		systemContent,
	)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	fmt.Println("Response:")
	debugPrint(result)

	fmt.Println("Responce only with text")
	for _, part := range result.Candidates[0].Content.Parts {
		fmt.Println(part.Text)
	}
}

func debugPrint[T any](r *T) {
	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response))
}
