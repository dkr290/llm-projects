package models

import (
	"encoding/json"
	"strings"
)

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Payload struct {
	Messages         []Message      `json:"messages"`
	Model            string         `json:"model"`
	FrequencyPenalty int            `json:"frequency_penalty"`
	MaxTokens        int            `json:"max_tokens"`
	PresencePenalty  int            `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	Stop             any            `json:"stop"`
	Stream           bool           `json:"stream"`
	StreamOptions    any            `json:"stream_options"`
	Temperature      int            `json:"temperature"`
	TopP             int            `json:"top_p"`
	Tools            any            `json:"tools"`
	ToolChoice       string         `json:"tool_choice"`
	Logprobs         bool           `json:"logprobs"`
	TopLogprobs      any            `json:"top_logprobs"`
}

func CreatePayload(systemPrompt, userMessage string) *strings.Reader {
	payload := Payload{
		Messages: []Message{
			{
				Content: systemPrompt,
				Role:    "system",
			},
			{
				Content: userMessage, // Now using a variable
				Role:    "user",
			},
		},
		Model:            "deepseek-chat",
		FrequencyPenalty: 0,
		MaxTokens:        2048,
		PresencePenalty:  0,
		ResponseFormat:   ResponseFormat{Type: "text"},
		Stop:             nil,
		Stream:           false,
		StreamOptions:    nil,
		Temperature:      1,
		TopP:             1,
		Tools:            nil,
		ToolChoice:       "none",
		Logprobs:         false,
		TopLogprobs:      nil,
	}

	jsonPayload, _ := json.Marshal(payload)
	return strings.NewReader(string(jsonPayload))
}
