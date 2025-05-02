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
	models    []string
}

func NewHandler(debugFlag bool, apiKey string, models []string) (*AMLController, error) {
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
		models:    models,
	}, nil
}

func (ac *AMLController) SearchHandler(
	c fuego.ContextWithBody[models.SearchRequest],
) (models.SearchResponse, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	body, err := c.Body()
	if err != nil {
		return models.SearchResponse{}, fuego.BadRequestError{
			Title: "invalid request payload",
			Err:   err,
		}
	}
	targetName := body.TargetName
	if targetName == "" {
		return models.SearchResponse{}, fuego.BadRequestError{
			Title: "target name is empty",
			Err:   err,
		}
	}
	ac.logger.Info().Str("targetName", targetName).Msg("Starting AML check")

	responseText, err := ac.ValidateLanguage(body.Language, targetName, ctx, "", "", true)
	if err != nil {
		if err.Error() == "not supported language" {
			return models.SearchResponse{}, fuego.BadRequestError{
				Title: "Language validation error choose supported language",
				Err:   err,
			}
		} else {
			return models.SearchResponse{}, fuego.BadRequestError{
				Title:  "Unknown Error",
				Detail: err.Error(),
			}
		}
	}
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
