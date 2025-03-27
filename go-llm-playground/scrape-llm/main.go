package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"screpe-llm/pkg/models"
	"screpe-llm/pkg/scraper"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/joho/godotenv"
)

type Scraper struct {
	Title string
	Text  []string
}

func main() {
	title, text := scraper.ScraperWeb("edwarddonner.com")

	systemPrompt := `You are an assistant that analyzes the contents of a website
                   and provides a short summary, ignoring text that might be navigation related.
                  R espond in markdown.`

	userprmpt := UserPrompt(title, text)

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API_KEY from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not found in .env file")
	}
	Chat(systemPrompt, userprmpt, apiKey)
}

func Chat(systemPrompt, userPrompt, apiKey string) {
	url := "https://api.deepseek.com/chat/completions"
	method := "POST"

	payload := models.CreatePayload(systemPrompt, userPrompt)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	mdfile := markdown.ToHTML(body, p, renderer)
	// Wrap in a full HTML structure
	// Save to an HTML file
	file, err := os.Create("output.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(mdfile))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("HTML file created successfully: output.html")
}

func UserPrompt(title string, text string) string {
	userPrompt := fmt.Sprintf("You are looking at a website titled %s", title)

	userPrompt += `\nThe contents of this website is as follows;
	                please provide a short summary of this website in markdown.
                  If it includes news or announcements, then summarize these too.\n\n`

	userPrompt += text

	return userPrompt
}
