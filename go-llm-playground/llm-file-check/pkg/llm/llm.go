package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

func LLmreplace(
	file []byte,
	ollamaHost string,
	model string,
) {
	prompt := fmt.Sprintf(`

You are an assistant that categorizes and sorts grocery items.

Here is a list of grocery items:

%s

Please:

1. Categorize these items into appropriate categories such as Produce, Dairy, Meat, Bakery, Beverages, etc.
2. Sort the items alphabetically within each category.
3. Present the categorized list in a clear and organized manner, using bullet points or numbering.

`, string(file))

	fmt.Println(prompt)
	os.Setenv("OLLAMA_HOST", ollamaHost) // Default value
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()
	genReq := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	// Capture the full response
	var fullResponse string
	genResp := func(r api.GenerateResponse) error {
		if r.Response == "" {
			return fmt.Errorf("empty responce from ollama")
		}
		fullResponse += r.Response
		return nil
	}
	err = client.Generate(ctx, genReq, genResp)
	if err != nil {
		log.Fatal("Error generating response from Ollama:", err)
	}
	// Write modified YAML to a new file
	err = os.WriteFile("categorized_grocery_list.txt", []byte(fullResponse), 0644)
	if err != nil {
		log.Fatal("Error writing categorized_grocery_list:", err)
	}

	fmt.Println("✅ Updated grocery_list.txt saved as categorized_grocery_list.txt!")
}
