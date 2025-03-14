package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {
	// Parse the URL string into a *url.URL
	apiURL, err := url.Parse("http://172.22.0.4/ollama")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	client := api.NewClient(apiURL, &http.Client{})
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel() // Always call cancel to release resources

	promt := "Write and introduction for an article about the future of AI."
	err = generareText(ctx, client, promt)
	if err != nil {
		fmt.Println(err)
	}
}

func generareText(
	ctx context.Context,
	client *api.Client,
	prompt string,
) error {
	fullPrompt := fmt.Sprintf(
		"Correct any spelling and grammer mistakes in the following text: %s",
		prompt,
	)

	req := &api.GenerateRequest{
		Model:  "deepseek-r1:1.5b",
		Prompt: fullPrompt,

		// set streaming to false
		Stream: new(bool),
	}
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println("AI generated content")
		fmt.Println(resp.Response)
		return nil
	}
	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
