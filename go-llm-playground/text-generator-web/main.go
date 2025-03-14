package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text-generator-web/app"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
		return errors.New("need ollama url")
	}
	fmt.Println("using the model", *model)
	// Parse the URL string into a *url.URL
	apiURL, err := url.Parse(*ollamaURL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %v", err)
	}

	client := api.NewClient(apiURL, &http.Client{})
	h := app.New(client, *model)

	engine := html.New("./views", ".html")
	api := fiber.New(fiber.Config{
		Views: engine,
	})
	api.Post("/generate", h.GenerateText)
	api.Get("/", h.IndexHandler)
	return api.Listen(port)
}

func getEnvs() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "0.0.0.0:3000"
	}
}
