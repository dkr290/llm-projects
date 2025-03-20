package main

import (
	"fmt"
	"log"
	"log/slog"
	"pdf-extract-opensearch/pkg/collection"
	"pdf-extract-opensearch/pkg/pdf"

	"github.com/spf13/cobra"
)

var (
	initVectors bool

	// userQuery string
	// chatModel   = "llama3.2"
	// embedModel  = "nomic-embed-text"
	// ollamaUrl   = "http://172.22.0.5/ollama"
	opensearch  string
	username    string
	password    string
	vectorIndex string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mycli",
		Short: "A CLI tool with OpenSearch and vector embeddings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opensearch == "" {
				return fmt.Errorf("error: --host is required")
			}
			if username == "" {
				return fmt.Errorf("error: --user is required")
			}
			if password == "" {
				return fmt.Errorf("error: --pass is required")
			}
			fmt.Println("Running CLI tool...")
			return nil
		},
	}

	// Define flags
	rootCmd.PersistentFlags().BoolVar(&initVectors, "init", false, "Initialize vector embeddings")
	rootCmd.PersistentFlags().StringVar(&opensearch, "host", "", "OpenSearch host")
	rootCmd.PersistentFlags().StringVar(&vectorIndex, "index", "vector_store", "Vector store index")
	rootCmd.PersistentFlags().StringVar(&username, "user", "", "OpenSearch username")
	rootCmd.PersistentFlags().StringVar(&password, "pass", "", "OpenSearch password")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	//	defer cancel()
	if initVectors {
		embed_run()
	}

	// if userQuery != "" {
	// 	r := response.New(chatModel, embedModel, ollamaUrl, qUrl, vectorCollection, userQuery, ctx)
	// 	if err := r.QuestionResponse(); err != nil {
	// 		slog.Error("unexpected error of question responce", "error", err)
	// 	}
	//
	// }
}

func embed_run() {
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
	// mdTextChunks, err := p.AddMetadata(docs, "BOI US FinCEN")
	_, err = p.AddMetadata(docs, "BOI US FinCEN")
	if err != nil {
		slog.Error("error metadata ", "error", err)
		return
	}
	fmt.Println("Adding metadata to the chunks")
	err = collection.CreateIndex(vectorIndex, opensearch, username, password)
	if err != nil {
		slog.Error("error creating opensearhc index", "error", err)
	}

	// err = p.GenEmbeddings(
	// 	mdTextChunks,
	// 	"nomic-embed-text",
	// 	"http://172.22.0.5/ollama",
	// 	"http://qdrant.172.22.0.5.nip.io:6333",
	// 	"vector_store1",
	// )
	// if err != nil {
	// 	slog.Error("error adding the embeddings", "error", err)
	// }
	// fmt.Println("Finish with vector embedding")
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
