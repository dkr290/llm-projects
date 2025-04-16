package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"press-detective/internal/logging"
	"press-detective/internal/models"
	"press-detective/internal/utils"

	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"

	"google.golang.org/genai"
)

type AMLController struct {
	client    *genai.Client
	logger    zerolog.Logger
	debugFlag bool
}

func NewAMLController(debugFlag bool) (*AMLController, error) {
	logger := logging.NewContextLogger("AMLController")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Error().Msg("GEMINI_API_KEY environment variable not set")
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

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

	result, err := ac.client.Models.GenerateContent(ctx, "gemini-2.0-flash",
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
