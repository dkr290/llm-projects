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

	"github.com/dkr290/llm-projects/go-llm-playground/customer-support-chat/go-supchat-llm/app"
	"github.com/dkr290/llm-projects/go-llm-playground/customer-support-chat/go-supchat-llm/models"
	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

var port string

func main() {
	getEnvs()
	log.Fatal(Run())
}

var FAQ_DB = map[string]string{
	"order tracking":           "You can track your order by logging into your account and navigating to 'My Orders'.",
	"return policy":            "We accept returns within 30 days. Visit our Returns page to initiate a return.",
	"customer support contact": "You can reach customer support at support@example.com or call us at +1-800-555-1234.",
	"payment methods":          "We accept Visa, MasterCard, PayPal, and Apple Pay for secure transactions.",
	"shipping details":         "Orders are processed within 24 hours. Standard shipping takes 3-5 business days.",
}

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
		DB:             FAQ_DB,
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
