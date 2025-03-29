package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Title string
	Text  []string
	Links []string
}

func ScraperWeb(website string) (string, string, string) {
	c := colly.NewCollector(
		colly.AllowedDomains(website),
		colly.MaxDepth(1),
	)

	var tt Scraper
	scriptRegex := regexp.MustCompile(
		`(?i)(function\s*\(|var\s+\w+|window\.\w+|document\.\w+|parentElement|insertBefore|_stq\.push|JSON\.parse|classList\.add)`,
	)
	jsonRegex := regexp.MustCompile(`(?i)\{.*?[:].*?\}`) // Matches JSON-like text
	longLineRegex := regexp.MustCompile(`[{}()\[\];]+`)  // Detects code-like lines

	// Get page title

	c.OnHTML("title", func(h *colly.HTMLElement) {
		tt.Title = h.Text
	})

	// Find and print all links
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content := e.Text
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !scriptRegex.MatchString(trimmed) &&
				!jsonRegex.MatchString(trimmed) &&
				!longLineRegex.MatchString(trimmed) {
				tt.Text = append(tt.Text, trimmed+"\n")
			}
		}
	})

	// Extract all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			tt.Links = append(tt.Links, link)
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
	allHref := strings.Join(tt.Links, ",")
	return tt.Title, allLines, allHref
}
