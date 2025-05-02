package handlers

import (
	"context"
	"errors"
	"fmt"
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
	var (
		result  *genai.GenerateContentResponse
		err     error
		lastErr error
	)

	for i, modelName := range ac.models {
		ac.logger.Info().Msg("Attempting to use model " + modelName)
		result, err = ac.client.Models.GenerateContent(ctx, modelName,
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
			ac.logger.Error().Err(err).Msg("Gemini API call failed with model " + modelName)
			lastErr = err
			if i == len(ac.models)-1 {
				return "", fmt.Errorf("failed to call all models with Gemini API: %w", lastErr)
			}
			continue
		} else {
			return utils.ExtractResponseText(result), nil
		}

	}
	return "", lastErr
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
	standardized_lang := strings.ToLower(lang)

	switch standardized_lang {

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
