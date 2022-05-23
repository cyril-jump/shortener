package storage

type DB interface {
	GetBaseURL(userID, shortURL string) (string, error)
	GetAllURLsByUserID(userID string) ([]ModelURL, error)
	SetShortURL(userID, shortURL, baseURL string) error
	Close()
}

type Users interface {
	GetUserID(userName string) (string, error)
	CreateCookie(userID string) (string, error)
	CheckCookie(cookieOld, userID string) bool
	SetUserID(userName string)
}

type ModelURL struct {
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"original_url"`
}
