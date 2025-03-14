package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dkr290/go-llm-playground/text-corrector/go-text-llm/helpers"
	"github.com/dkr290/go-llm-playground/text-corrector/go-text-llm/models"

	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

type Handlers struct {
	client         *api.Client
	conf           models.Config
	ContextTimeout time.Duration
}

func New(api *api.Client, c models.Config) *Handlers {
	return &Handlers{
		client: api,
		conf:   c,
	}
}

func (h *Handlers) GenerateText(c *fiber.Ctx) error {
	var req models.JsonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Validate required fields
	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both 'prompt' is required",
		})
	}
	promt := fmt.Sprintf("%s\n\n%s", h.conf.ConfigPrompt, req.Prompt)
	ctx, cancel := context.WithTimeout(context.Background(), h.conf.ContextTimeout)
	defer cancel()

	err := ensureModelExists(ctx, h.client, h.conf.Model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to ensure model exists: %v", err))
	}

	generatedText, err := helpers.GenerareText(ctx, h.client, promt, h.conf.Model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Cannot generate the text: " + err.Error()})
	}
	return c.JSON(fiber.Map{
		"generated_text": generatedText,
	})
}

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
