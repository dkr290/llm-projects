package main

import (
	"flag"
	"fmt"
	"llm-file-check/pkg/config"
	"llm-file-check/pkg/llm"
	"log/slog"
	"os"
)

func main() {
	// Define flags
	model := flag.String("model", "", "The model to use")
	serverURL := flag.String("server", "", "The server URL")
	file := flag.String("file", "", "The values file")
	help := flag.Bool("help", false, "Show usage information")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Parse the flags
	flag.Parse()
	if *help {
		showUsage()
	}

	if *model == "" || *serverURL == "" || *file == "" {
		showUsage()
	}

	config := config.NewConfig(*file)
	f, err := config.ReadHelmValues()
	if err != nil {
		logger.Error("Error read helm values", "error", err)
	}
	llm.LLmreplace(f, *serverURL, *model)
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  -model model to serve")
	fmt.Println("  -server the http or https address to ollama server")
	fmt.Println("  -file the helm values file")

	os.Exit(0)
}
