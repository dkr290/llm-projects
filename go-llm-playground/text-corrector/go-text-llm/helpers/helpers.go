package helpers

import (
	"context"
	"fmt"

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
		return "", err
	}

	return generatedText, nil
}

// pullModel pulls the specified model using the Ollama client
func PullModel(ctx context.Context, client *api.Client, model string) error {
	req := &api.PullRequest{
		Name: model,
	}

	respFunc := func(resp api.ProgressResponse) error {
		fmt.Printf("Pulling model: %s and status %s\n", resp.Digest, resp.Status)
		return nil
	}

	err := client.Pull(ctx, req, respFunc)
	if err != nil {
		return fmt.Errorf("failed to pull model: %v", err)
	}

	return nil
}
