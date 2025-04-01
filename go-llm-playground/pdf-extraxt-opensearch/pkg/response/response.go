package response

import (
	"context"
	"fmt"
	"strings"

	ops "github.com/opensearch-project/opensearch-go"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/vectorstores/opensearch"
)

type Questions struct {
	ChatModel        string
	CTX              context.Context
	EmbedModel       string
	OllamaUrl        string
	OpensearchClient *ops.Client
}

func New(
	chatModel, embedModel, ollamaUrl string,
	ctx context.Context, o *ops.Client,
) *Questions {
	return &Questions{
		ChatModel:        chatModel,
		EmbedModel:       embedModel,
		OllamaUrl:        ollamaUrl,
		OpensearchClient: o,
		CTX:              ctx,
	}
}

func (q Questions) QuestionResponse(UserQuery string) (map[string]any, error) {
	// Load chat model (Ollama 3.2)
	chatModel, err := ollama.New(
		ollama.WithModel(q.ChatModel), // Change to your chat model
		ollama.WithServerURL(q.OllamaUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading chat model: %v", err)
	}

	// Load embedding model (nomic-embed-text)
	embedModel, err := ollama.New(
		ollama.WithModel(q.EmbedModel),
		ollama.WithServerURL(q.OllamaUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading embedding model: %v", err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(embedModel)
	if err != nil {
		return nil, fmt.Errorf("new embedder error %v", err)
	}

	store, err := opensearch.New(q.OpensearchClient, opensearch.WithEmbedder(ollamaEmbeder))
	if err != nil {
		return nil, fmt.Errorf("error initializing Opensearch store %v", err)
	}

	// Example user query
	//	userQuery := "by when should I file if my business was established in 2013?"

	// Retrieve documents from Qdrant using a single query
	docs, err := store.SimilaritySearch(
		q.CTX,
		UserQuery,
		5,
	) // Retrieve top 5 results
	if err != nil {
		return nil, fmt.Errorf("error retrieving documents %v", err)
	}
	fmt.Println(docs)

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

	fmt.Println(finalChain)
	answer, err := finalChain.Call(q.CTX, map[string]any{
		"context":  contextText,
		"question": UserQuery,
	}, chains.WithMaxTokens(500), chains.WithMaxLength(10000))
	// chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 	if ctx.Err() != nil || len(chunk) == 0 {
	// 		return nil
	// 	}
	// 	fmt.Printf("%s", chunk)
	// 	return nil
	// }))
	if err != nil {
		return nil, fmt.Errorf("Error generating answer: %v", err)
	}
	fmt.Println(answer)
	// Print AI-generated answer
	return answer, nil
}
