package handlers

import (
	"gemini-press-api/internal/models"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	h, _ := NewHandler(false, "test", []string{"test"})

	ctx := fuego.NewMockContextNoBody()
	response, err := h.RootHandler(ctx)
	assert.NoError(t, err)
	assert.Equal(t, models.MessageResponse{Message: "Hello, from go press detector app!"}, response)
}
