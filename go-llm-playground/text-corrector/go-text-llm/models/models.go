package models

import "time"

type JsonRequest struct {
	Prompt string `json:"prompt"`
}

type Config struct {
	Model          string
	ConfigPrompt   string
	ContextTimeout time.Duration
}
