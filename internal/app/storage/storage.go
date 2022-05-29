package storage

type DB interface {
	GetBaseURL(shortURL string) (string, error)
	GetAllURLsByUserID(userID string) ([]ModelURL, error)
	SetShortURL(userID, shortURL, baseURL string) error
	Ping() error
	Close() error
}

type Users interface {
	CreateCookie(userID string) (string, error)
	CheckCookie(tokenString string) (string, bool)
}

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
