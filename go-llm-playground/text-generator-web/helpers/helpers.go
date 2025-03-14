package helpers

import (
	"context"
	"fmt"
	"regexp"

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
		generatedText = formatAsHTML(resp.Response)
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

// formatAsHTML formats plain text as HTML
func formatAsHTML(text string) string {
	// Example: Format code blocks, headings, and lists
	// Replace Markdown-like syntax with HTML tags
	formattedText := text

	// Format headings (e.g., "# Heading" -> <h1>Heading</h1>)
	formattedText = regexp.MustCompile(`(?m)^#\s+(.*)$`).
		ReplaceAllString(formattedText, `<h1>$1</h1>`)
	formattedText = regexp.MustCompile(`(?m)^##\s+(.*)$`).
		ReplaceAllString(formattedText, `<h2>$1</h2>`)
	formattedText = regexp.MustCompile(`(?m)^###\s+(.*)$`).
		ReplaceAllString(formattedText, `<h3>$1</h3>`)

	// Format code blocks (e.g., ```go ... ``` -> <pre><code class="go">...</code></pre>)
	formattedText = regexp.MustCompile("(?s)```go\\s*(.*?)```").
		ReplaceAllString(formattedText, `<pre><code class="go">$1</code></pre>`)
	formattedText = regexp.MustCompile("(?s)```\\s*(.*?)```").
		ReplaceAllString(formattedText, `<pre><code>$1</code></pre>`)

	// Format inline code (e.g., `code` -> <code>code</code>)
	formattedText = regexp.MustCompile("`(.*?)`").ReplaceAllString(formattedText, `<code>$1</code>`)

	// Format lists (e.g., "- item" -> <li>item</li>)
	formattedText = regexp.MustCompile(`(?m)^-\s+(.*)$`).
		ReplaceAllString(formattedText, `<li>$1</li>`)
	formattedText = regexp.MustCompile(`(?m)^\*\s+(.*)$`).
		ReplaceAllString(formattedText, `<li>$1</li>`)

	// Wrap lists in <ul> tags
	formattedText = regexp.MustCompile(`(?s)(<li>.*</li>)`).
		ReplaceAllString(formattedText, `<ul>$1</ul>`)

	return formattedText
}
