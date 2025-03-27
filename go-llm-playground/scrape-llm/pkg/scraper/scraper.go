package scraper

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

func ScraperWeb(website string) (string, string) {
	c := colly.NewCollector(
		colly.AllowedDomains(website),
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
	err := c.Visit("https://" + website)
	if err != nil {
		log.Fatal("error Visit ", err)
	}

	allLines := strings.Join(tt.Text, "\n") // Joins with newline between lines
	return tt.Title, allLines
}
