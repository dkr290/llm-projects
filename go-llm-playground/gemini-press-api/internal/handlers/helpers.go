// Package handlers contains the HTTP handlers for the application.
package handlers

import (
	"context"
	"errors"
	"fmt"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/utils"
	"strings"
	"time"

	"google.golang.org/genai"
)

func (ac *AMLController) makeGeminiRequest(
	ctx context.Context,
	targetName string,
	prompt string,
) (string, error) {
	var lastErr error
	for keyIndex, key := range ac.apiKeys {
		ac.logger.Info().Msgf("Attempting with API key %d/%d", keyIndex+1, len(ac.apiKeys))
		client, err := createGeminiClient(key)
		if err != nil {
			ac.logger.Error().
				Err(err).
				Msgf("Failed to create client with API key %d/%d", keyIndex+1, len(ac.apiKeys))
			lastErr = err
			continue
		}
		for i, modelName := range ac.models {
			ac.logger.Info().
				Msgf("Attempting model %s (%d/%d) with API key %d/%d", modelName, i+1, len(ac.models), keyIndex+1, len(ac.apiKeys))

				// using separate function for gemini to be able to test
			result, err := generateGeminiResponce(ctx, client, modelName, targetName, prompt)

			if err != nil {
				ac.logger.Error().
					Err(err).
					Msgf("Failed: model %s with API key %d/%d", modelName, keyIndex+1, len(ac.apiKeys))
				lastErr = err
				continue
			} else {
				ac.logger.Info().Msgf("Success: model %s with API key %d/%d", modelName, keyIndex+1, len(ac.apiKeys))
				return utils.ExtractResponseText(result), nil
			}

		}
		// Only return error if this was the last key
		if keyIndex == len(ac.apiKeys)-1 {
			return "", fmt.Errorf("failed to call all models with all API keys: %w", lastErr)
		}

	}
	return "", lastErr
}

func createGeminiClient(key string) (*genai.Client, error) {
	logger := logging.NewContextLogger("AMLController")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Gemini client")
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return client, err
}

func ValidateDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format: must be YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format: must be YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, fmt.Errorf(
			"invalid date range: start date must be before or equal to end date",
		)
	}

	return startDate, endDate, nil
}

func (ac *AMLController) ValidateLanguage(
	lang string, targetName string,
	ctx context.Context, startDate, endDate string, isGeneric bool,
) (responseText string, err error) {
	standardizedLang := strings.ToLower(lang)

	switch standardizedLang {

	case "english":
		responseText, err = ac.makeGeminiRequest(
			ctx,
			targetName,
			utils.GetEnglishPrompt(startDate, endDate, isGeneric),
		)
		if err != nil {
			return "", err
		}

	case "french":
		responseText, err = ac.makeGeminiRequest(
			ctx,
			targetName,
			utils.GetFrenchPrompt(startDate, endDate, isGeneric),
		)
		if err != nil {
			return "", err
		}

	default:
		return "", errors.New("not supported language")
	}
	return responseText, nil
}
