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
	// Variables to track progress
	var totalSize int64
	var downloaded int64
	var currentDigest string

	respFunc := func(resp api.ProgressResponse) error {
		// Initialize total size if provided in response
		if resp.Total > 0 && totalSize == 0 {
			totalSize = resp.Total
		}

		// Update current digest if changed
		if resp.Digest != "" && resp.Digest != currentDigest {
			currentDigest = resp.Digest
			downloaded = 0 // Reset downloaded bytes for new layer
			fmt.Printf("\nPulling layer: %s\n", currentDigest)
		}

		// Update downloaded bytes
		if resp.Completed > 0 {
			downloaded = resp.Completed
		}
		// Calculate and display progress
		if totalSize > 0 {
			percentage := float64(downloaded) / float64(totalSize) * 100
			fmt.Printf("\rStatus: %s | Progress: %.2f%% (%d/%d bytes)",
				resp.Status,
				percentage,
				downloaded,
				totalSize,
			)
		} else {
			// Fallback when total size isn't available
			fmt.Printf("\rStatus: %s | Downloaded: %d bytes",
				resp.Status,
				downloaded,
			)
		}
		return nil
	}
	err := client.Pull(ctx, req, respFunc)
	if err != nil {
		return fmt.Errorf("failed to pull model: %v", err)
	}

	fmt.Println("\nModel pull completed successfully")
	return nil
}
