package handlers

import (
	"fmt"
	"pdf-extract-opensearch/models"
	"pdf-extract-opensearch/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	ResponceConfig response.Questions
}

func NewHandlers(r response.Questions) *Handlers {
	return &Handlers{
		ResponceConfig: r,
	}
}

func (h *Handlers) GenerateText(c *fiber.Ctx) error {
	var req models.JsonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request " + err.Error()})
	}
	// Validate required fields
	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "'prompt' is required",
		})
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	var resp map[string]any
	var err error

	if resp, err = h.ResponceConfig.QuestionResponseNew(req.Prompt); err != nil {
		fmt.Println(resp["text"])
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Cannot generate response: " + err.Error()})
	}
	responseText, ok := resp["text"].(string)
	fmt.Println(resp["text"])

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid response format from AI",
		})
	}

	return c.JSON(fiber.Map{
		"AI Responce:": responseText,
	})
}
