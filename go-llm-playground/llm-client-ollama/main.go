package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

// Helper function to create a *bool from a bool value
func boolPtr(b bool) *bool {
	return &b
}

func main() {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		os.Setenv("OLLAMA_HOST", "http://172.22.0.3/ollama") // Default value
	}
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()
	// resp, err := client.List(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(resp)

	// messages := []api.Message{
	// 	{
	// 		Role:    "user",
	// 		Content: "why the sky is blue?",
	// 	},
	// }
	//
	// req := &api.ChatRequest{
	// 	Model:    "llama3.2",
	// 	Messages: messages,
	// 	Stream:   boolPtr(true),
	// }
	//
	// respFunc := func(resp api.ChatResponse) error {
	// 	fmt.Print(resp.Message.Content)
	// 	return nil
	// }
	//
	// err = client.Chat(ctx, req, respFunc)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("")

	// api.GenerateRequest{
	// 	Model:  "llama3.2",
	// 	Prompt: "why is the sky blue ?",
	// }
	// showReq := api.ShowRequest{
	// 	Model: "llama3.2",
	// }
	// showResp, err := client.Show(ctx, &showReq)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(showResp.ModelInfo)

	req := &api.CreateRequest{
		Model:  "knowitall",
		From:   "llama3.2",
		System: "You are very smart assistant who knows everything about oceans.",
		Parameters: map[string]any{
			"temperature": 0.3,
		},
	}
	// Function to track progress
	progressFunc := func(update api.ProgressResponse) error {
		fmt.Println("Progress:", update.Status)
		return nil
	}
	err = client.Create(ctx, req, progressFunc)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	genReq := &api.GenerateRequest{
		Model:  "knowitall",
		Prompt: "why is the ocean is salty ?",
	}

	genResp := func(r api.GenerateResponse) error {
		fmt.Print(r.Response)
		return nil
	}
	err = client.Generate(ctx, genReq, genResp)
	if err != nil {
		log.Fatal(err)
	}

	delReq := &api.DeleteRequest{
		Model: "knowitall",
	}
	err = client.Delete(ctx, delReq)
	if err != nil {
		log.Fatal(err)
	}
}
