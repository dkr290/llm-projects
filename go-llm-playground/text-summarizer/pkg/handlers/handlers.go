package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"text-summarizer/pkg/summarizer"

	"github.com/ollama/ollama/api"
)

type Handlers struct {
	template     *template.Template
	ollamaClient *api.Client
	summarizer   *summarizer.Config
}

func New(s *summarizer.Config, olClient *api.Client, temp *template.Template) *Handlers {
	return &Handlers{
		template:     temp,
		ollamaClient: olClient,
		summarizer:   s,
	}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	err := h.template.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) SummarizeHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed")
	}

	text := r.FormValue("text")
	if text == "" {
		return fmt.Errorf("error: Please input text to summarize")
	}

	summary, err := h.summarizer.SummarizeText(text)
	if err != nil {
		return err
	}
	fmt.Println(summary)

	// Instead of rendering a full template, send a partial HTML response for HTMX
	w.Header().Set("Content-Type", "text/html")
	_, err = fmt.Fprintf(
		w,
		`<div class="result"><h3>Summary:</h3><p>%s</p></div>`,
		template.HTMLEscapeString(summary),
	)
	if err != nil {
		return fmt.Errorf("failed to write response: %v", err)
	}
	return nil
}
