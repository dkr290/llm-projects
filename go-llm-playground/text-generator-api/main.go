package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text-generator-api/app"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

var port string

func main() {
	getEnvs()
	log.Fatal(Run())
}

func Run() error {
	ollamaURL := flag.String("ollamaurl", "", "ollama url")
	model := flag.String("model", "deepseek-r1:1.5b", "ollama model to use")

	flag.Parse()
	if *ollamaURL == "" {
		return errors.New("need ollama model")
	}
	fmt.Println("using the model", *model)
	// Parse the URL string into a *url.URL
	apiURL, err := url.Parse(*ollamaURL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %v", err)
	}

	client := api.NewClient(apiURL, &http.Client{})
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel() // Always call cancel to release resources
	h := app.New(ctx, client, *model)

	api := fiber.New()
	api.Post("/generate", h.GenerateText)
	return api.Listen(port)
}

func getEnvs() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = ":3000"
	}
}
