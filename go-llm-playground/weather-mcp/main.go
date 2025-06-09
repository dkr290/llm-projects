package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// WeatherData represents the weather model
type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type WeatherService struct{}

func (s *WeatherService) GetWeather(
	ctx context.Context,
	req *protocol.CallToolRequest,
) (*protocol.CallToolResult, error) {
	weather := WeatherData{
		Temperature: rand.Float64()*40 - 10,
		Condition:   randomCondition(),
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: weather,
			},
		},
	}
}

func main() {
	// Create SSE transport server
	transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to create transport server: %v", err)
	}

	// Initialize MCP server
	mcpServer, err := server.NewServer(transportServer)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	tool, err := protocol.NewTool(
		"weather",
		"Get current weather temperature and condition",
		WeatherData{},
	)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
		return
	}

	// Start server
	if err = mcpServer.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func randomCondition() string {
	conditions := []string{"Sunny", "Rainy", "Cloudy", "Stormy", "Snowy"}
	return conditions[rand.Intn(len(conditions))]
}
