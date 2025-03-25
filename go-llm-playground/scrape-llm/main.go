package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Title string
	Text  []string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("edwarddonner.com"),
	)

	var tt Scraper

	c.OnHTML("title", func(h *colly.HTMLElement) {
		tt.Title = h.Text
	})

	// Find and print all links
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content := e.Text
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				tt.Text = append(tt.Text, trimmed)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	err := c.Visit("https://edwarddonner.com/")
	if err != nil {
		log.Fatal(err)
	}

	// Print results
	fmt.Println("\nTitle:")
	fmt.Println(tt.Title)

	fmt.Println("\nText Content:")
	for _, line := range tt.Text {
		fmt.Println(line)
	}
}
