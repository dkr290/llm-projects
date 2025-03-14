package llmquery

import (
	"context"
	"fmt"

	ops "github.com/opensearch-project/opensearch-go"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/opensearch"
)

func QuestionResponse(ctx context.Context, model string, url string, client *ops.Client) error {
	// Initialize the Ollama client with the provided model and server URL
	ollamaLLM, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(url),
	)
	if err != nil {
		return fmt.Errorf("ollama: %w", err)
	}
	// Format a prompt to direct the model what to do with the content and
	// the question.
	prompt := `Use the following pieces of information to answer the user's question.
	If you don't know the answer, say that you don't know.
	
	Question: %s

	Answer the question and provide additional helpful information, but be concise.

	Responses should be properly formatted to be easily read.
	`

	question := `what is the document about?`

	finalPrompt := fmt.Sprintf(prompt, question)
	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		return fmt.Errorf("new embedder error %v", err)
	}

	store, err := opensearch.New(client,
		opensearch.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		return fmt.Errorf("error opensearch with embedder %v", err)
	}

	OptionVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80),
	}
	retrreaver := vectorstores.ToRetriever(store, 1, OptionVector...)
	resDocs, err := retrreaver.GetRelevantDocuments(ctx, finalPrompt)
	if err != nil {
		return fmt.Errorf("failed to retreive documents %v", err)
	}

	fmt.Println(resDocs)
	return nil
}
