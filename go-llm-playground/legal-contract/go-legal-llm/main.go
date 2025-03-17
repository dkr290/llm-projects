package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dkr290/llm-projects/go-llm-playground/legal-contract/go-legal-llm/app"
	"github.com/dkr290/llm-projects/go-llm-playground/legal-contract/go-legal-llm/models"
	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

var port string

func main() {
	getEnvs()
	log.Fatal(Run())
}

// Legal Document Templates

func Run() error {
	ollamaURL := flag.String("ollamaurl", "", "ollama url")
	model := flag.String("model", "deepseek-r1:1.5b", "ollama model to use")
	language := flag.String("language", "English", "default to english language")

	contextTimeout := flag.Duration(
		"timeout",
		30*time.Minute,
		"the context timeout like 2 * time.Minute",
	)

	flag.Parse()
	if *ollamaURL == "" {
		return errors.New("need ollama url")
	}
	fmt.Println("using the model", *model)
	// Parse the URL string into a *url.URL
	apiURL, err := url.Parse(*ollamaURL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %v", err)
	}
	client := api.NewClient(apiURL, &http.Client{})
	h := app.New(client, models.Config{
		Model:          *model,
		Lang:           *language,
		ContextTimeout: *contextTimeout,
	})
	api := fiber.New()
	api.Post("/generate", h.GenerateText)
	return api.Listen(port)
}

func getEnvs() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "0.0.0.0:3000"
	}
}
