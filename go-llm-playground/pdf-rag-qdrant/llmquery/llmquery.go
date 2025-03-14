package llmquery

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

func QuestionResponse(ctx context.Context) error {
	// Load chat model (Ollama 3.2)
	chatModel, err := ollama.New(
		ollama.WithModel("llama3.2"), // Change to your chat model
		ollama.WithServerURL("http://172.22.0.5/ollama"),
	)
	if err != nil {
		log.Fatal("Error loading chat model:", err)
	}

	// Load embedding model (nomic-embed-text)
	embedModel, err := ollama.New(
		ollama.WithModel("nomic-embed-text"),
		ollama.WithServerURL("http://172.22.0.5/ollama"),
	)
	if err != nil {
		log.Fatal("Error loading embedding model:", err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(embedModel)
	if err != nil {
		return fmt.Errorf("new embedder error %v", err)
	}
	// Create a new Qdrant vector store.
	url, err := url.Parse("http://qdrant.172.22.0.5.nip.io/")
	if err != nil {
		log.Fatal(err)
	}
	store, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithCollectionName("vector_store"),
		qdrant.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		return fmt.Errorf("error initializing Qdrant store %v", err)
	}

	// Example user query
	userQuery := "goroutines short example"

	// Retrieve documents from Qdrant using a single query
	docs, err := store.SimilaritySearch(ctx, userQuery, 3) // Retrieve top 3 results
	if err != nil {
		return fmt.Errorf("error retrieving documents %v", err)
	}

	// Combine retrieved documents as context
	var retrievedDocs []string
	for _, doc := range docs {
		retrievedDocs = append(retrievedDocs, doc.PageContent)
	}
	contextText := strings.Join(retrievedDocs, "\n")
	// Define the template with placeholders
	template := `Answer the question based ONLY on the following context:
Context:
{{.context}}
User Question: {{.question}}
Answer:`
	// Specify the input variables used in the template
	inputVariables := []string{"context", "question"}
	// Define chat prompt with retrieved context
	chatPrompt := prompts.NewPromptTemplate(template, inputVariables)
	// Generate final response using chat model
	finalChain := chains.NewLLMChain(chatModel, chatPrompt)

	answer, err := finalChain.Call(ctx, map[string]any{
		"context":  contextText,
		"question": userQuery,
	}, chains.WithMaxTokens(500), chains.WithMaxLength(10000))
	// chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 	if ctx.Err() != nil || len(chunk) == 0 {
	// 		return nil
	// 	}
	// 	fmt.Printf("%s", chunk)
	// 	return nil
	// }))
	if err != nil {
		log.Fatal("Error generating answer:", err)
	}
	fmt.Println("")
	// Print AI-generated answer
	fmt.Println("AI Response:\n", answer["text"])

	return nil
}
