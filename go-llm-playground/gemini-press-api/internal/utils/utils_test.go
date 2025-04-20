package utils

import (
	"gemini-press-api/internal/models"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestParseGeminiResponse(t *testing.T) {
	logger := zerolog.Nop() // Create a no-op logger for testing

	tests := []struct {
		name      string
		inputText string
		debugFlag bool
		want      models.SearchResponse
	}{
		{
			name:      "Red flag found, with link and summary",
			inputText: "red_flag_found: TRUE\nlinks:\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx\n summary: This is a test summary.",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "red_flag_found: TRUE\nlinks:\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx\n summary: This is a test summary.",
				RedFlagFound:   true,
				Links: []string{
					"https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq",
					"https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx",
				},
				Summary: "This is a test summary.",
			},
		},
		{
			name:      "Red flag not found",
			inputText: "red_flag_found: FALSE\nlinks:\nsummary:\n",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "red_flag_found: FALSE\nlinks:\nsummary:\n",
				RedFlagFound:   false,
				Links:          nil,
				Summary:        "",
			},
		},
		{
			name:      "Red flag found, no links, with summary",
			inputText: "red_flag_found: true\nsummary: Another summary here.",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "red_flag_found: true\nsummary: Another summary here.",
				RedFlagFound:   true,
				Links:          nil,
				Summary:        "Another summary here.",
			},
		},
		{
			name:      "Red flag found, with links, no summary",
			inputText: "red_flag_found: TRUE\nlinks:\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx\n summary: ",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "red_flag_found: TRUE\nlinks:\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq\n- https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx\n summary: ",
				RedFlagFound:   true,
				Links: []string{
					"https://vertexaisearch.cloud.google.com/grounding-api-redirect/AWQVq",
					"https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxx",
				},
				Summary: "",
			},
		},
		{
			name:      "No red flag found keyword",
			inputText: "Some text without the red flag.",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "Some text without the red flag.",
				RedFlagFound:   false,
				Links:          nil,
				Summary:        "",
			},
		},
		{
			name:      "Red flag found, multiple links in one line",
			inputText: "red_flag_found: true This has links: https://one.com and https://two.org and maybe ftp://not-http.com.",
			debugFlag: false,
			want: models.SearchResponse{
				RawLLMResponse: "red_flag_found: true This has links: https://one.com and https://two.org and maybe ftp://not-http.com.",
				RedFlagFound:   true,
				Links: []string{
					"https://one.com",
					"https://two.org",
				}, // ftp is not matched by https?
				Summary: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseGeminiResponse(tt.inputText, logger, tt.debugFlag)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGeminiResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
