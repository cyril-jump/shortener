package dto

//Data Transfer Object Packet

type ModelURL struct {
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"original_url"`
}

type ModelURLBatchRequest struct {
	CorID   string `json:"correlation_id"`
	BaseURL string `json:"original_url"`
}

type ModelURLBatchResponse struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type ModelRequestURL struct {
	BaseURL string `json:"url"`
}

type ModelResponseURL struct {
	ShortURL string `json:"result"`
}
