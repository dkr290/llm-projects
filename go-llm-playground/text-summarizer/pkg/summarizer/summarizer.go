package summarizer

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type Config struct {
	client *api.Client
	model  string
	ctx    context.Context
}

func New(client *api.Client, model string, ctx context.Context) *Config {
	return &Config{
		model:  model,
		client: client,
		ctx:    ctx,
	}
}

func toPtr[T any](t T) *T {
	return &t
}

func (c *Config) SummarizeText(text string) (string, error) {
	var summary string

	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: fmt.Sprintf("Summarize the following text in **3 bullet points**:\n\n%s", text),
		Stream: toPtr(false),
	}
	respFunc := func(resp api.GenerateResponse) error {
		summary = resp.Response
		return nil
	}
	err := c.client.Generate(c.ctx, req, respFunc)
	if err != nil {
		return "", fmt.Errorf("cannot query ollama %v", err)
	}

	return summary, nil
}
