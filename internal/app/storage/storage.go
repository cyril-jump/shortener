package storage

type DB interface {
	GetBaseURL(shortURL string) (string, error)
	GetAllURLsByUserID(userID string) ([]ModelURL, error)
	SetShortURL(userID, shortURL, baseURL string) error
	Close()
}

type Users interface {
	CreateCookie(userID string) (string, error)
	CheckCookie(tokenString string) (string, bool)
}

type ModelURL struct {
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"original_url"`
}
