package storage

import "github.com/cyril-jump/shortener/internal/app/dto"

//DB interface
type DB interface {
	GetBaseURL(shortURL string) (string, error)
	GetAllURLsByUserID(userID string) ([]dto.ModelURL, error)
	SetShortURL(userID, shortURL, baseURL string) error
	DelBatchShortURLs(tasks []dto.Task) error
	Ping() error
	Close() error
}

//Users interface
type Users interface {
	CreateToken(userID string) (string, error)
	CheckToken(tokenString string) (string, bool)
}

//Cfg Config interface
type Cfg interface {
	Get(key string) (string, error)
}

//InWorker interface
type InWorker interface {
	Do(t dto.Task)
	Loop() error
}
