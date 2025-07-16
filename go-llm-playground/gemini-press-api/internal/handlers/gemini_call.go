package handlers

import (
	"context"

	"google.golang.org/genai"
)

func generateGeminiResponce(
	ctx context.Context,
	client *genai.Client,
	modelName string,
	targetName, prompt string,
) (*genai.GenerateContentResponse, error) {
	result, err := client.Models.GenerateContent(ctx, modelName,
		genai.Text(targetName),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: prompt}},
			},
			Tools: []*genai.Tool{{
				GoogleSearch: &genai.GoogleSearch{},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
