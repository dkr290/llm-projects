package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OllamaURL is the endpoint for the Ollama API
const OllamaURL = "http://172.22.0.4/ollama/api/generate"

// RequestPayload defines the structure of the request to Ollama
type RequestPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ResponsePayload defines the structure of the response from Ollama
type ResponsePayload struct {
	Response string `json:"response"`
}

func summarizeText(text string) string {
	// Create the payload
	payload := RequestPayload{
		Model:  "llama3.2", // Replace with your desired model
		Prompt: fmt.Sprintf("Summarize the following text in **3 bullet points**:\n\n%s", text),
		Stream: false,
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("Error: Failed to marshal payload - %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send POST request to Ollama API
	resp, err := client.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error: Failed to connect to Ollama server - %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Error: %s", resp.Status)
	}

	// Decode the response
	var result ResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Sprintf("Error: Failed to decode response - %v", err)
	}

	// Return the summary or a fallback message
	if result.Response == "" {
		return "No summary generated."
	}
	return result.Response
}

func main() {
	// Example usage
	sampleText := `
	Artificial intelligence (AI) is transforming industries worldwide. From healthcare to finance,
	AI systems are improving efficiency, reducing costs, and enabling new capabilities. In
	healthcare, AI assists with diagnostics and personalized treatment plans. In finance, it
	enhances fraud detection and algorithmic trading. However, challenges like ethical concerns,
	data privacy, and job displacement remain significant hurdles to its widespread adoption.
	`

	summary := summarizeText(sampleText)
	fmt.Println(summary)
}
