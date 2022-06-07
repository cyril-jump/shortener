package storage

import "github.com/cyril-jump/shortener/internal/app/dto"

type DB interface {
	GetBaseURL(shortURL string) (string, error)
	GetAllURLsByUserID(userID string) ([]dto.ModelURL, error)
	SetShortURL(userID, shortURL, baseURL string) error
	DelBatchShortURLs(tasks []dto.Task)
	Ping() error
	Close() error
}

type Users interface {
	CreateToken(userID string) (string, error)
	CheckToken(tokenString string) (string, bool)
}

type Cfg interface {
	Get(key string) (string, error)
}

type InWorker interface {
	Do(t *dto.Task)
	Loop() error
}
