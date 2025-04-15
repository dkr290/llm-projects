package models

// SearchRequest defines the structure for the search request payload.
type SearchRequest struct {
	TargetName string `json:"target_name" example:"Osama Bin Laden" description:"Name to check"`
}

// SearchResponse defines the structure for the search response payload.
type SearchResponse struct {
	RedFlagFound   bool     `json:"red_flag_found"             example:"true"`
	RawLLMResponse *string  `json:"raw_llm_response,omitempty" example:"Based on the search results..."   description:"Raw response from the LLM"`
	Links          []string `json:"links,omitempty"            example:"https://example.com/news"         description:"List of relevant links"`
	Summary        *string  `json:"summary,omitempty"          example:"Found multiple fraud allegations" description:"Summary of the findings"`
}

// HealthCheck defines the structure for the health check response.
type HealthCheck struct {
	Status            string `json:"status"             example:"OK"`
	GeminiConfigured  bool   `json:"gemini_configured"  example:"true"`
	GeminiOperational bool   `json:"gemini_operational" example:"true"`
}
