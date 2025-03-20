package main

import (
	"flag"
	"fmt"
	"log/slog"
	"pdf-extract-opensearch/pkg/collection"
	"pdf-extract-opensearch/pkg/pdf"
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
	flag.BoolVar(&initVectors, "init", false, "init set to truth is to init the vectgor embedding")
	flag.StringVar(&opensearch, "host", "", "opensearch host")
	flag.StringVar(&vectorIndex, "index", "vector_store", "the vector store index")
	// flag.StringVar(&username, "user", "", "the username")
	// flag.StringVar(&password, "pass", "", "password for opensearch")

	flag.Parse()

	fmt.Println(opensearch)
	fmt.Println(initVectors)
	fmt.Println(vectorIndex)

	// if username == "" || password == "" || opensearch == "" {
	// 	fmt.Println("Some parameters are empty")
	// 	os.Exit(1)
	// }
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
