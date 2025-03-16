package models

import "time"

type JsonRequest struct {
	Prompt string `json:"prompt"`
}

type Config struct {
	Model          string
	DB             map[string]string
	ContextTimeout time.Duration
	Lang           string
}
