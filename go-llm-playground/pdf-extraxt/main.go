package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"pdf-extract/pkg/collection"
	"pdf-extract/pkg/pdf"
	"pdf-extract/pkg/response"
	"time"
)

var (
	initVectors bool

	userQuery        string
	chatModel        = "llama3.2"
	embedModel       = "nomic-embed-text"
	ollamaUrl        = "http://172.22.0.5/ollama"
	qUrl             = "http://qdrant.172.22.0.5.nip.io:6333"
	vectorCollection = "vector_store1"
)

func main() {
	flag.BoolVar(&initVectors, "init", false, "init set to truth is to init the vectgor embedding")
	flag.StringVar(&userQuery, "query", "", "asking question")

	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	if initVectors {
		embed_run()
	}

	if userQuery != "" {
		r := response.New(chatModel, embedModel, ollamaUrl, qUrl, vectorCollection, userQuery, ctx)
		if err := r.QuestionResponse(); err != nil {
			slog.Error("unexpected error of question responce", "error", err)
		}

	}
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
	mdTextChunks, err := p.AddMetadata(docs, "BOI US FinCEN")
	if err != nil {
		slog.Error("error metadata ", "error", err)
		return
	}
	fmt.Println("Adding metadata to the chunks")

	err = collection.CreateCollection("vector_store1", "qdrant-grpc.172.22.0.5.nip.io")
	if err != nil {
		slog.Error("error creating collection", "error", err)
		return
	}
	fmt.Println("Creating the qdrant collection")

	err = p.GenEmbeddings(
		mdTextChunks,
		"nomic-embed-text",
		"http://172.22.0.5/ollama",
		"http://qdrant.172.22.0.5.nip.io:6333",
		"vector_store1",
	)
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
