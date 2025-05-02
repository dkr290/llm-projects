package models

type MessageResponse struct {
	Message string `json:"message"`
}

// SearchRequest defines the structure for the search request payload.
type SearchRequest struct {
	TargetName string `json:"target_name"`
	Language   string `json:"language"    example:"english" description:"Input language english, french"`
}

// SearchResponse defines the structure for the search response payload.
type SearchResponse struct {
	RedFlagFound   bool     `json:"red_flag_found"             example:true                               description:"TRUE or FALSE"`
	RawLLMResponse string   `json:"raw_llm_response,omitempty" example:"Based on the search results..."   description:"Raw response from the LLM"`
	Links          []string `json:"links,omitempty"            example:["https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxxx","https://https://vertexaisearch.cloud.google.com/grounding-api-redirect/xxxx"]         description:"List of relevant links"`
	Summary        string   `json:"summary,omitempty"          example:"Found multiple fraud allegations" description:"Summary of the findings"`
}
