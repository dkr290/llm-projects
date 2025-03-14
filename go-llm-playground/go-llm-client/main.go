package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	// Define flags
	model := flag.String("model", "llama3.2", "The model to use")
	serverURL := flag.String("server", "http://ollama", "The server URL")
	query := flag.String("query", "what is golang", "The query to ask the model")
	timeTaken := flag.Bool(
		"time-taken",
		false,
		"true or false to see how much time in seconds took for the request",
	)
	help := flag.Bool("help", false, "Show usage information")

	// Parse the flags
	flag.Parse()
	if *help {
		showUsage()
	}

	// Initialize the Ollama client with the provided model and server URL
	llm, err := ollama.New(
		ollama.WithModel(*model),
		ollama.WithServerURL(*serverURL),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // Always call cancel to release resources
	completion, timeDuration, err := GenerateFromPrompt(ctx, llm, *query, *timeTaken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:\n", completion)
	if timeDuration != 0 {
		fmt.Println("Time taken in seconds", timeDuration)
	}
}

func GenerateFromPrompt(
	ctx context.Context,
	llm *ollama.LLM,
	query string,
	enabledTimeTaken bool,
) (string, time.Duration, error) {
	var elapsedTime time.Duration
	var completion string
	var err error
	if enabledTimeTaken {
		startTime := time.Now()

		completion, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
		if err != nil {
			return "", elapsedTime, err
		}
		elapsedTime = time.Since(startTime)
		return completion, elapsedTime, err
	}

	completion, err = llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		return "", elapsedTime, err
	}

	return completion, elapsedTime, nil
}

// showUsage prints usage information and exits
func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  -model model to serve")
	fmt.Println("  -server the http or https address to ollama server")
	fmt.Println("  -query the query to ask")
	fmt.Println("  -time-taken true default if false")
	os.Exit(0)
}
