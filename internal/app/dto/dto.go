package dto

//Data Transfer Object Packet

// ModelURL struct
type ModelURL struct {
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"original_url"`
}

// ModelURLBatchRequest struct
type ModelURLBatchRequest struct {
	CorID   string `json:"correlation_id"`
	BaseURL string `json:"original_url"`
}

// ModelURLBatchResponse struct
type ModelURLBatchResponse struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// ModelRequestURL struct
type ModelRequestURL struct {
	BaseURL string `json:"url"`
}

// ModelResponseURL struct
type ModelResponseURL struct {
	ShortURL string `json:"result"`
}

// Task struct
type Task struct {
	ID       string
	ShortURL string
}

// Stat struct
type Stat struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
