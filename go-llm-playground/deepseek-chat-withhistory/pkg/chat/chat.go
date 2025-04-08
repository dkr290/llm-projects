package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type ChatScruct struct {
	SystemMessage string
	HumanMessage  string
	ChatHistory   []llms.MessageContent
}

func New(systemMessage string) *ChatScruct {
	return &ChatScruct{
		SystemMessage: systemMessage,
		ChatHistory: []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, systemMessage),
		},
	}
}

func (c *ChatScruct) DsChat(ctx context.Context) error {
	// Initialize the OpenAI client with Deepseek model
	llm, err := openai.New(
		openai.WithModel("deepseek-chat"),
	)
	if err != nil {
		return fmt.Errorf("error connectint to deepseek %v", err)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Chat with AI. Type 'exit' to quit.")

	for {
		fmt.Print("You: ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)

		if strings.ToLower(query) == "exit" {
			break
		}

		humanMsg := llms.TextParts(llms.ChatMessageTypeHuman, query)
		c.ChatHistory = append(c.ChatHistory, humanMsg)

		// Generate content with streaming to see both reasoning and final answer in real-time
		completion, err := llm.GenerateContent(
			ctx,
			c.ChatHistory,
			llms.WithMaxTokens(2000),
			llms.WithTemperature(0.7),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("error generate context %v", err)
		} // Add AI response to history

		// Access the reasoning content and final answer separately
		if len(completion.Choices) > 0 {
			aiResponse := completion.Choices[0].Content
			aiMsg := llms.TextParts(llms.ChatMessageTypeAI, aiResponse)
			c.ChatHistory = append(c.ChatHistory, aiMsg)
			fmt.Printf("\nAI: %s\n", aiResponse)
		}

	}
	fmt.Println("\n----Message History----")
	for _, msg := range c.ChatHistory {
		switch msg.Role {
		case llms.ChatMessageTypeSystem:
			fmt.Printf("System: %s\n", msg.Parts)
		case llms.ChatMessageTypeHuman:
			fmt.Printf("Human: %s\n", msg.Parts)
		case llms.ChatMessageTypeAI:
			fmt.Printf("AI: %s\n", msg.Parts)
		}
	}
	return nil
}
