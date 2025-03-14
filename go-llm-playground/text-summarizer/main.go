package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"text-summarizer/pkg/handlers"
	"text-summarizer/pkg/helpers"
	"text-summarizer/pkg/summarizer"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {
	run()
}

func run() {
	// Parse the URL string into a *url.URL
	apiURL, err := url.Parse("http://172.22.0.4/ollama")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	client := api.NewClient(apiURL, &http.Client{})
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel() // Always call cancel to release resources
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
	// prompt := `Summmarize the following text:
	//  Artificial intelligence is transforming the industried acress the world. AI models like DeppSeek R1 enable
	//  businesses to automate tasks,  analize large datasets, and enahance productiviti. With advancements of AI, applications
	//  range from virtusal assistants to predictive analytics and personalized reccomendations.
	//  `
	summarizer := summarizer.New(client, "llama3.2", ctx)
	h := handlers.New(summarizer, client, temp)

	http.HandleFunc("/", helpers.MakeHandler(h.HomeHandler))
	http.HandleFunc("/summarize", helpers.MakeHandler(h.SummarizeHandler))

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
