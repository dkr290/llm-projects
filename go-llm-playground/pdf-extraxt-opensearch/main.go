package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"pdf-extract-opensearch/pkg/collection"
	"pdf-extract-opensearch/pkg/pdf"
	"time"

	"github.com/opensearch-project/opensearch-go"

	"github.com/spf13/cobra"
)

var (
	initVectors bool

	// userQuery string
	chatModel   string
	embedModel  string
	ollamaUrl   string
	opsHost     string
	username    string
	password    string
	vectorIndex string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mycli",
		Short: "A CLI tool with OpenSearch and vector embeddings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opsHost == "" {
				return fmt.Errorf("error: --host is required")
			}
			if username == "" {
				return fmt.Errorf("error: --user is required")
			}
			if password == "" {
				return fmt.Errorf("error: --pass is required")
			}
			if ollamaUrl == "" {
				return fmt.Errorf("error: --ollamaurl is required")
			}
			fmt.Println("Running CLI tool...")
			return nil
		},
	}

	// Define flags
	rootCmd.PersistentFlags().BoolVar(&initVectors, "init", false, "Initialize vector embeddings")
	rootCmd.PersistentFlags().StringVar(&opsHost, "opshost", "", "OpenSearch host")
	rootCmd.PersistentFlags().StringVar(&vectorIndex, "index", "vector_store", "Vector store index")
	rootCmd.PersistentFlags().StringVar(&username, "opsuser", "", "OpenSearch username")
	rootCmd.PersistentFlags().StringVar(&password, "opspass", "", "OpenSearch password")
	rootCmd.PersistentFlags().
		StringVar(&chatModel, "chatmodel", "llama3.2", "the chat models to use")
	rootCmd.PersistentFlags().
		StringVar(&embedModel, "embedmodel", "nomic-embed-text", "the embedding model to use")
	rootCmd.PersistentFlags().StringVar(&ollamaUrl, "ollamaurl", "", "the ollama url")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	// Create a client instance
	cfg := opensearch.Config{
		Addresses: []string{
			opsHost, // Replace with your OpenSearch URL

		},
		Username: username,
		Password: password,
	}
	client, err := opensearch.NewClient(cfg)
	if err != nil {
		slog.Error("error creating the client: ", "error", err)
		os.Exit(1)
	}

	if initVectors {
		embed_run(client, embedModel)
	}

	models := []string{embedModel, chatModel}

	for _, m := range models {
		err := collection.EnsureModelExists(ctx, ollamaUrl, m)
		if err != nil {
			log.Fatal("error with downloading models", err)
		}
	}

	// if userQuery != "" {
	// 	r := response.New(chatModel, embedModel, ollamaUrl, qUrl, vectorCollection, userQuery, ctx)
	// 	if err := r.QuestionResponse(); err != nil {
	// 		slog.Error("unexpected error of question responce", "error", err)
	// 	}
	//
	// }
}

func embed_run(client *opensearch.Client, embedmodel string) {
	// new PDF reading
	p := pdf.New("./BOI.pdf")
	if err := p.ReadPdf(); err != nil {
		slog.Error("error reading pdf ", "error", err)
		return
	}
	fmt.Println("Reading the pdf")
	// split text in chunks
	docs, err := p.SplitText()
	if err != nil {
		slog.Error("error splitting the text", "error", err)
		return
	}
	fmt.Println("Splitting the documents to chunks")

	// add metadata to the chunks
	mdTextChunks, err := p.AddMetadata(docs, "BOI US FinCEN")
	if err != nil {
		slog.Error("error metadata ", "error", err)
		return
	}

	fmt.Println("Adding metadata to the chunks")
	err = collection.CreateIndex(vectorIndex, client)
	if err != nil {
		slog.Error("error creating opensearhc index", "error", err)
	}
	err = p.GenEmbeddings(mdTextChunks, embedModel, ollamaUrl, client, vectorIndex)
	if err != nil {
		slog.Error("error adding the embeddings", "error", err)
	}
	fmt.Println("Finish with vector embedding")
}

// curl -X PUT "http://localhost:6333/collections/my_collection" \
//      -H "Content-Type: application/json" \
//      -d '{
//            "vectors": {
//                "size": 128,
//                "distance": "Cosine"
//            },
//            "shard_number": 6,
//            "replication_factor": 3
//          }'
