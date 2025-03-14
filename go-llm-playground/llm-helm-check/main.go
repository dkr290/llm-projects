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
	valuesFile := flag.String("values", "", "The values file")
	chartFile := flag.String("chart", "", "The chart file")
	internalRegisty := flag.String("registry", "", "The internal registry to replace")
	help := flag.Bool("help", false, "Show usage information")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Parse the flags
	flag.Parse()
	if *help {
		showUsage()
	}

	if *model == "" || *serverURL == "" || *valuesFile == "" || *chartFile == "" ||
		*internalRegisty == "" {
		showUsage()
	}

	config := config.NewConfig(*valuesFile, *chartFile)
	values, chart, err := config.ReadHelmValues()
	if err != nil {
		logger.Error("Error read helm values", "error", err)
	}
	llm.LLmreplace(values, chart, *internalRegisty, *serverURL, *model)
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  -model model to serve")
	fmt.Println("  -server the http or https address to ollama server")
	fmt.Println("  -values the helm values file")
	fmt.Println("  -chart the chart file")
	fmt.Println("  -registry the internal registry")

	os.Exit(0)
}
