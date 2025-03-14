package helpers

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

func GenerareText(
	ctx context.Context,
	client *api.Client,
	prompt string,
	model string,
) (string, error) {
	req := &api.GenerateRequest{
		// Model:  "deepseek-r1:1.5b",
		Model:  model,
		Prompt: prompt,

		// set streaming to false
		Stream: new(bool),
	}
	var generatedText string
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		generatedText = resp.Response
		return nil
	}
	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	return generatedText, nil
}
