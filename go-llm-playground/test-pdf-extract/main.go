package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func readPDF(filePath string) string {
	outDir := "tmp"                 // Specify the output directory
	selectedPages := []string{"1-"} // Extract all pages
	conf := model.NewDefaultConfiguration()

	err := api.ExtractContentFile(filePath, outDir, selectedPages, conf)
	if err != nil {
		log.Fatal(err)
	}

	// Read the extracted content from the output directory
	// Assuming the content is extracted to a text file
	content, err := os.ReadFile(outDir + "/extracted.txt")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

// This loop iterates over the text in steps of chunkSize.
// i starts at 0 and increments by chunkSize in each iteration.

func splitIntoChunks(text string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}
	return chunks
}

func main() {
	filePath := "BOI.pdf"
	content := readPDF(filePath)
	chunks := splitIntoChunks(content, 1000) // Adjust chunk size as needed
	fmt.Println(chunks)
}
