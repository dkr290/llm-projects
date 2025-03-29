package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"marketing-broshure/pkg/models"
	"marketing-broshure/pkg/scraper"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_, _, links := scraper.ScraperWeb("edwarddonner.com")

	systemPrompt := `You are provided with a list of links found on a webpage.
You are able to decide which of the links would be most relevant to include in brochure about the company,
such as links to About page, or Company page, or Career/Jobs pages.
You Should respond in JSON as in this example.
	  {
	    "links": [
	      { "type:"about page","url":"https://full.url/goes/here/about" },
	      { "type": "careers page": "url": "https://another.full.url/careers"}
	     ]
	  }
	`
	userprmpt := UserPrompt("https://edwarddonner.com", links)

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
	parseLinks(body)
}

func UserPrompt(websiteUrl, links string) string {
	userPrompt := fmt.Sprintf("Here is the list of links on the website of  %s\n", websiteUrl)

	userPrompt += `Please decide which of theese are relevant web links for the brochure about the company, respond with the full https URL,
	Do not include Terms of Service , Privacy email links.`

	userPrompt += "\nLinks (some might be relative links):\n"
	userPrompt += fmt.Sprintf("\n %s", links)
	return userPrompt
}

func parseLinks(jsonData []byte) {
	// Define a struct to match the JSON structure
	type Link struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	fmt.Println(string(jsonData))
	type Data struct {
		Links []Link `json:"links"`
	}
	// Parse JSON
	var data Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	for _, link := range data.Links {
		fmt.Printf("Type: %s, URL: %s\n", link.Type, link.URL)
	}
}
