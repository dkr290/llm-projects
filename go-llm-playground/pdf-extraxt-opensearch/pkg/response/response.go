package response

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ops "github.com/opensearch-project/opensearch-go"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

type Questions struct {
	ChatModel        string
	CTX              context.Context
	EmbedModel       string
	OllamaUrl        string
	OpensearchClient *ops.Client
	VectorIndex      string
}

func New(
	chatModel, embedModel, ollamaUrl string,
	ctx context.Context, o *ops.Client, vIndex string,
) *Questions {
	return &Questions{
		ChatModel:        chatModel,
		EmbedModel:       embedModel,
		OllamaUrl:        ollamaUrl,
		OpensearchClient: o,
		CTX:              ctx,
		VectorIndex:      vIndex,
	}
}

// func (q Questions) QuestionResponse(UserQuery string) (map[string]any, error) {
// 	// Load chat model (Ollama 3.2)
// 	chatModel, err := ollama.New(
// 		ollama.WithModel(q.ChatModel), // Change to your chat model
// 		ollama.WithServerURL(q.OllamaUrl),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error loading chat model: %v", err)
// 	}
//
// 	// Load embedding model (nomic-embed-text)
// 	embedModel, err := ollama.New(
// 		ollama.WithModel(q.EmbedModel),
// 		ollama.WithServerURL(q.OllamaUrl),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error loading embedding model: %v", err)
// 	}
// 	ollamaEmbeder, err := embeddings.NewEmbedder(embedModel)
// 	if err != nil {
// 		return nil, fmt.Errorf("new embedder error %v", err)
// 	}
//
// 	store, err := opensearch.New(q.OpensearchClient, opensearch.WithEmbedder(ollamaEmbeder))
// 	if err != nil {
// 		return nil, fmt.Errorf("error initializing Opensearch store %v", err)
// 	}
//
// 	// Example user query
// 	//	userQuery := "by when should I file if my business was established in 2013?"
//
// 	// Retrieve documents from Qdrant using a single query
// 	docs, err := store.SimilaritySearch(
// 		q.CTX,
// 		UserQuery,
// 		5,
// 	) // Retrieve top 5 results
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving documents %v", err)
// 	}
// 	fmt.Println(docs)
//
// 	// Combine retrieved documents as context
// 	var retrievedDocs []string
// 	for _, doc := range docs {
// 		retrievedDocs = append(retrievedDocs, doc.PageContent)
// 	}
// 	contextText := strings.Join(retrievedDocs, "\n")
// 	// Define the template with placeholders
// 	template := `Answer the question based ONLY on the following context:
// 	Context:
// 	{{.context}}
// 	User Question: {{.question}}
// 	Answer:`
// 	// Specify the input variables used in the template
// 	inputVariables := []string{"context", "question"}
// 	// Define chat prompt with retrieved context
// 	chatPrompt := prompts.NewPromptTemplate(template, inputVariables)
// 	// Generate final response using chat model
// 	finalChain := chains.NewLLMChain(chatModel, chatPrompt)
//
// 	fmt.Println(finalChain)
// 	answer, err := finalChain.Call(q.CTX, map[string]any{
// 		"context":  contextText,
// 		"question": UserQuery,
// 	}, chains.WithMaxTokens(500), chains.WithMaxLength(10000))
// 	// chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
// 	// 	if ctx.Err() != nil || len(chunk) == 0 {
// 	// 		return nil
// 	// 	}
// 	// 	fmt.Printf("%s", chunk)
// 	// 	return nil
// 	// }))
// 	if err != nil {
// 		return nil, fmt.Errorf("Error generating answer: %v", err)
// 	}
// 	fmt.Println(answer)
// 	// Print AI-generated answer
// 	return answer, nil
// }

func (q Questions) QuestionResponseNew(UserQuery string) (map[string]any, error) {
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

	queryEmbedding, err := ollamaEmbeder.EmbedQuery(q.CTX, UserQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	vectorJSON, err := json.Marshal(queryEmbedding)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vector: %w", err)
	}

	// Build manual k-NN query
	query := fmt.Sprintf(`{
        "size": 5,
        "query": {
            "knn": {
                "embedding": {
                    "vector": %s,
                    "k": 5
                }
            }
        }
    }`, vectorJSON)

	// Execute search directly through OpenSearch client
	res, err := q.OpensearchClient.Search(
		q.OpensearchClient.Search.WithIndex(q.VectorIndex),
		q.OpensearchClient.Search.WithBody(strings.NewReader(query)),
		q.OpensearchClient.Search.WithContext(q.CTX),
		q.OpensearchClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("opensearch search error: %w", err)
	}
	defer res.Body.Close()

	// Parse results
	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	var retrievedDocs []string
	if hits, ok := result["hits"].(map[string]any); ok {
		if hitList, ok := hits["hits"].([]any); ok {
			for _, hit := range hitList {
				if source, ok := hit.(map[string]any)["_source"].(map[string]any); ok {
					if content, ok := source["content"].(string); ok {
						retrievedDocs = append(retrievedDocs, content)
					}
				}
			}
		}
	}

	if len(retrievedDocs) == 0 {
		return nil, fmt.Errorf("no documents found for query")
	}

	// Combine context with safe joining
	contextText := strings.Join(retrievedDocs, "\n---\n")
	// Create prompt template
	template := `Answer the question based ONLY on the following context:
Context:
{{.context}}

User Question: {{.question}}
Answer in complete sentences and be as helpful as possible:`

	inputVariables := []string{"context", "question"}
	chatPrompt := prompts.NewPromptTemplate(template, inputVariables)

	// Create chain and execute
	finalChain := chains.NewLLMChain(chatModel, chatPrompt)
	answer, err := finalChain.Call(q.CTX, map[string]any{
		"context":  contextText,
		"question": UserQuery,
	}, chains.WithMaxTokens(500), chains.WithMaxLength(1000))
	if err != nil {
		return nil, fmt.Errorf("error generating answer: %w", err)
	}
	// Return formatted response
	return map[string]any{
		"answer":    answer,
		"documents": retrievedDocs, // Optional: include source documents
	}, nil
}
