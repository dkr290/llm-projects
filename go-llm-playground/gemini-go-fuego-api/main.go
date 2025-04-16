package main

import (
	"fmt"
	"gemini-go-fuego-api/internal/services"
	"gemini-go-fuego-api/internal/utils"
	"log"
	"os"
	"strings"
)

func main() {
	logger := utils.NewLogger()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the GEMINI_API_KEY environment variable")
		return
	}
	adverseMediaService, err := services.NewService(apiKey)
	if err != nil {

		logger.Error().
			Str("Error:", err.Error()).
			Msg("Error initializing Adverse Media Service")

		return
	}
	resp, err := adverseMediaService.GenerateReport("Who is Someone")
	if err != nil {
		if strings.Compare(resp.Summary, "") == 0 {
			fmt.Println("No response from llm")
			return
		}

		logger.Error().
			Str("Error:", err.Error()).
			Msg("Error initializing Adverse Media Service")

		return
	}

	fmt.Println("")
	fmt.Println("Summary", resp.Summary)

	fmt.Println("links", resp.Links)
	fmt.Println("RedFlags", resp.RedFlagFound)
	fmt.Println("Raw LLM responce", resp.RawLLMResponse)
}
