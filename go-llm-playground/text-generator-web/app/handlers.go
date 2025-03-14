package app

import (
	"context"
	"fmt"
	"log"
	"text-generator-web/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

type Handlers struct {
	client *api.Client
	model  string
}

func New(api *api.Client, model string) *Handlers {
	return &Handlers{
		client: api,
		model:  model,
	}
}

func (h *Handlers) GenerateText(c *fiber.Ctx) error {
	prmpt := c.FormValue("prompt")
	if prmpt == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Prompt is required")
	} // Validate required fields
	words := c.FormValue("words")
	if words == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Number of words is required")
	} // Validate required fields

	promt := fmt.Sprintf("Generate %s words:\n\n%s", words, prmpt)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	err := ensureModelExists(ctx, h.client, h.model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to ensure model exists: %v", err))
	}

	generatedText, err := helpers.GenerareText(ctx, h.client, promt, h.model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to generate text: %v", err))
	}

	// Return the generated text as the response
	// log.Println(generatedText)
	return c.SendString(generatedText)
}

func (h *Handlers) IndexHandler(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

// ensureModelExists checks if the model exists locally and pulls it if it doesn't
func ensureModelExists(ctx context.Context, client *api.Client, model string) error {
	// Try to generate a dummy prompt to check if the model exists
	_, err := helpers.GenerareText(ctx, client, "test prompt", model)
	if err == nil {
		// Model exists locally
		return nil
	}

	// If the error indicates that the model is missing, pull it

	log.Printf("Model %s not found locally. Pulling...\n", model)
	return helpers.PullModel(ctx, client, model)
	// Return other errors
}
