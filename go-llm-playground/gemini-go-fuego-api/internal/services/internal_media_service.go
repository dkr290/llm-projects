package services

import (
	"context"
	"fmt"
	"gemini-go-fuego-api/models"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/genai"
)

// Service interface defines the operations for adverse media checking.
type Service interface {
	GenerateReport(targetName string) (*models.SearchResponse, error)
}

// service implements the Service interface.
type service struct {
	geminiClient *genai.Client
}

// NewService creates a new adverse media service using the genai client.
func NewService(geminiAPIKey string) (Service, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %w", err)
	}
	return &service{geminiClient: client}, nil
}

// Define the structure for the Gemini API response (simplified for now)
// We will be using the genai library's response types.

// Prompt template
const promptTemplate = `You are an Anti Money Laundering expert analyst who specializes in doing adverse media checks.
Respond ONLY in the following EXACT format:

red_flag_found: [TRUE|FALSE]
links: (only if red_flag_found is TRUE)
- https://example1.com
- https://example2.com
summary: (only if red_flag_found is TRUE)
Brief summary of findings (2-3 sentences maximum)

Search template to use:
"{TARGET_NAME}" AND (Scam OR Convict OR Fraud OR charged OR Terror OR radical OR guilty OR forced labor OR slavery OR embezzlement OR Scandal OR Theft OR Forgery OR Jailed OR illegal OR Evasion OR drugs OR Abuse OR Misconduct OR Fine OR Sanctions OR Corruption)

Verify any potential red flags by visiting the links before reporting them.
`

// GenerateReport interacts with the Gemini API to check for adverse media using the genai library.
func (s *service) GenerateReport(targetName string) (*models.SearchResponse, error) {
	ctx := context.Background()
	model := s.geminiClient.GenerativeModel("gemini-pro") // Or another suitable model

	prompt := strings.ReplaceAll(promptTemplate, "{TARGET_NAME}", targetName)

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ToolConfig{
			FunctionCalling: false, // Explicitly disable function calling for this prompt
			GoogleSearch:    &genai.GoogleSearchConfig{},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error generating content: %w", err)
	}

	var rawLLMResponse string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			rawLLMResponse = string(textPart)
		}
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
