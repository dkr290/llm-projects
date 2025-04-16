package utils

import (
	"gemini-press-api/internal/models"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/genai"
)

func ExtractResponseText(result *genai.GenerateContentResponse) string {
	if result == nil || result.Text() == "" {
		return ""
	}
	return result.Text()
}

func ParseGeminiResponse(
	responseText string,
	logger zerolog.Logger,
	debugFlag bool,
) models.SearchResponse {
	response := models.SearchResponse{RawLLMResponse: responseText}

	redFlagMatch := regexp.MustCompile(`(?i)red_flag_found:\s*(true|false)`)
	if match := redFlagMatch.FindStringSubmatch(responseText); len(match) > 1 {
		response.RedFlagFound = strings.EqualFold(match[1], "true")
	}

	if response.RedFlagFound {
		linkRegex := regexp.MustCompile(`https?://[^\s]+`)
		response.Links = linkRegex.FindAllString(responseText, -1)

		if summaryMatch := regexp.MustCompile(`(?s)summary:\s*(.*?)(\n\w+:|$)`).FindStringSubmatch(responseText); len(
			summaryMatch,
		) > 1 {
			response.Summary = strings.TrimSpace(summaryMatch[1])
		}
	}
	if debugFlag {
		logger.Debug().Interface("response", response).Msg("Parsed AML response")
	}
	return response
}
