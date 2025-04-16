package handlers

import (
	"context"
	"fmt"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/models"
	"gemini-press-api/internal/utils"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
	"google.golang.org/genai"
)

type AMLController struct {
	client    *genai.Client
	logger    zerolog.Logger
	debugFlag bool
	model     string
}

func NewHandler(debugFlag bool, apiKey string, model string) (*AMLController, error) {
	logger := logging.NewContextLogger("AMLController")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Gemini client")
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &AMLController{
		client:    client,
		logger:    logger,
		debugFlag: debugFlag,
		model:     model,
	}, nil
}

func (ac *AMLController) SearchHandler(
	c fuego.ContextWithBody[models.SearchRequest],
) (models.SearchResponse, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	body, err := c.Body()
	if err != nil {
		ac.logger.Error().Err(err).Msg("Failed to parse payload body")
		return models.SearchResponse{}, fmt.Errorf("payload unpacking failed")
	}

	targetName := body.TargetName

	ac.logger.Info().Str("targetName", targetName).Msg("Starting AML check")

	result, err := ac.client.Models.GenerateContent(ctx, ac.model,
		genai.Text(targetName),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: utils.GetPromptTemplate()}},
			},
			Tools: []*genai.Tool{{
				GoogleSearch: &genai.GoogleSearch{},
			}},
		},
	)
	if err != nil {
		ac.logger.Error().Err(err).Msg("Gemini API call failed")
		return models.SearchResponse{}, fmt.Errorf("API call failed")
	}

	responseText := utils.ExtractResponseText(result)
	return utils.ParseGeminiResponse(responseText, ac.logger, ac.debugFlag), nil
}

func (ac *AMLController) RootHandler(c fuego.ContextNoBody) (models.MessageResponse, error) {
	ac.logger.Info().Msg("Handling root request")
	return models.MessageResponse{Message: "Hello, from go press detector app!"}, nil
}

func (ac AMLController) PingHandler(c fuego.ContextNoBody) (models.MessageResponse, error) {
	ac.logger.Info().Msg("Handling ping request")
	return models.MessageResponse{Message: "pong"}, nil
}
