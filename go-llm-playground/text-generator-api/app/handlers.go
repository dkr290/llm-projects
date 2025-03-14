package app

import (
	"context"
	"fmt"
	"strings"
	"text-generator-api/helpers"
	"text-generator-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
)

type Handlers struct {
	ctx    context.Context
	client *api.Client
	model  string
}

func New(ctc context.Context, api *api.Client, model string) *Handlers {
	return &Handlers{
		ctx:    ctc,
		client: api,
		model:  model,
	}
}

func (h *Handlers) GenerateText(c *fiber.Ctx) error {
	var req models.JsonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	// Validate required fields
	if req.Prompt == "" || req.Wordlimit == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both 'prompt' and 'wordlimit' are required",
		})
	}
	promt := fmt.Sprintf("Generate %s words:\n\n%s", req.Wordlimit, req.Prompt)
	generatedText, err := helpers.GenerareText(h.ctx, h.client, promt, h.model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Cannot generate the text: " + err.Error()})
	}
	// Format the generated text to include new lines
	formattedText := fmt.Sprintf("<think>\n%s", generatedText)

	// Replace double new lines with single new lines (if needed)
	formattedText = strings.ReplaceAll(formattedText, "\n\n", "\n")

	return c.JSON(fiber.Map{
		"generated_text": formattedText,
	})
}
