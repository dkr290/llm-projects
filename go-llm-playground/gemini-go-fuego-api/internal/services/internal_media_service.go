package services

import (
	"context"
	"fmt"
	"gemini-go-fuego-api/models"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

// Service interface defines the operations media checking.
type Service interface {
	GenerateReport(targetName string) (*models.SearchResponse, error)
}

// service implements the Service interface.
type service struct {
	geminiClient *genai.Client
}

// NewService creates a new adverse media service using the genai client aind accepts api key.
func NewService(geminiAPIKey string) (Service, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %w", err)
	}
	return &service{geminiClient: client}, nil
}

// GenerateReport interacts with the Gemini API to check for adverse media using the genai library.
// this is where the actual reposnce from gemini is preocessed
func (s *service) GenerateReport(targetName string) (*models.SearchResponse, error) {
	ctx := context.Background()
	model := "gemini-2.0-flash" // Or another suitable model

	prompt := strings.ReplaceAll(promptTemplate, "{TARGET_NAME}", targetName)

	userContent := []*genai.Part{
		{Text: prompt},
	}
	systemContent := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}

	result, err := s.geminiClient.Models.GenerateContent(
		ctx,
		model,
		[]*genai.Content{{Parts: userContent}},
		systemContent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}
	var rawLLMResponse string
	for _, part := range result.Candidates[0].Content.Parts {
		rawLLMResponse = part.Text
	}

	// Parse and validate the Gemini response
	parsedResponse, err := parseGeminiResponse(rawLLMResponse)
	if err != nil {
		return nil, err
	}

	if !validateGeminiResponse(rawLLMResponse) {
		return nil, fmt.Errorf("AI returned unexpected response format: %s", rawLLMResponse)
	}

	return parsedResponse, nil
}

// validateGeminiResponse checks if the Gemini response contains the expected structure.
func validateGeminiResponse(text string) bool {
	re := regexp.MustCompile(`red_flag_found:\s*(?i)(true|false)`)
	return re.MatchString(text)
}

// parseGeminiResponse parses the Gemini response into a structured SearchResponse object.
func parseGeminiResponse(response_text string) (*models.SearchResponse, error) {
	rawResponseText := response_text
	responseText := strings.ReplaceAll(strings.TrimSpace(response_text), "\r\n", "\n")

	var redFlagFound bool
	redFlagMatch := regexp.MustCompile(`red_flag_found:\s*(?i)(true|false)`).
		FindStringSubmatch(responseText)
	if len(redFlagMatch) > 1 {
		redFlagFound = strings.ToLower(redFlagMatch[1]) == "true"
	}

	var links []string
	var summary *string

	if redFlagFound {
		linksSectionMatch := regexp.MustCompile(`links:(?s)(.*?)(summary:|red_flag_found:|\z)(?i)`).
			FindStringSubmatch(responseText)
		if len(linksSectionMatch) > 1 {
			linkMatches := regexp.MustCompile(`https?://\S+`).
				FindAllString(strings.TrimSpace(linksSectionMatch[1]), -1)
			links = linkMatches
		}

		summaryMatch := regexp.MustCompile(`summary:\s*(?s)(.*?)(red_flag_found:|links:|\z)(?i)`).
			FindStringSubmatch(responseText)
		if len(summaryMatch) > 1 {
			trimmedSummary := strings.TrimSpace(summaryMatch[1])
			summary = &trimmedSummary
		}
	}

	return &models.SearchResponse{
		RedFlagFound:   redFlagFound,
		RawLLMResponse: &rawResponseText,
		Links:          links,
		Summary:        summary,
	}, nil
}
