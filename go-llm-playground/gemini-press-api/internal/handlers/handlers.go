package handlers

import (
	"context"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/models"
	"gemini-press-api/internal/utils"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
)

type AMLController struct {
	logger    zerolog.Logger
	debugFlag bool
	models    []string
	apiKeys   []string
}

func NewHandler(debugFlag bool, apiKeys []string, models []string) (*AMLController, error) {
	logger := logging.NewContextLogger("AMLController")

	return &AMLController{
		logger:    logger,
		debugFlag: debugFlag,
		models:    models,
		apiKeys:   apiKeys,
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
	ac.logger.Debug().Msg("Handling root request")
	return models.MessageResponse{Message: "Hello, from go press detector app!"}, nil
}

func (ac AMLController) PingHandler(c fuego.ContextNoBody) (models.MessageResponse, error) {
	ac.logger.Debug().Msg("Handling ping request")
	return models.MessageResponse{Message: "pong"}, nil
}
