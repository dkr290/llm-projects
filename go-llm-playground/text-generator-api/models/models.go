package models

type JsonRequest struct {
	Prompt    string `json:"prompt"`
	Wordlimit string `json:"wordlimit"`
}
