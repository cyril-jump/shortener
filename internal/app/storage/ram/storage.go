package ram

import (
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// DB

type DB struct {
	DataCache map[string][]storage.ModelURL
}

//constructor

func NewDB() *DB {

	return &DB{DataCache: make(map[string][]storage.ModelURL)}
}

func (D *DB) GetBaseURL(userID, shortURL string) (string, error) {
	if _, ok := D.DataCache[userID]; ok {
		for _, val := range D.DataCache[userID] {
			if val.ShortURL == shortURL {
				return val.BaseURL, nil
			}
		}
	}
	return "", errs.ErrNoContent
}

func (D *DB) GetAllURLsByUserID(userID string) ([]storage.ModelURL, error) {
	if _, ok := D.DataCache[userID]; ok {
		return D.DataCache[userID], nil
	}
	return nil, errs.ErrNoContent
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	modelURL := storage.ModelURL{
		ShortURL: shortURL,
		BaseURL:  baseURL,
	}
	if _, ok := D.DataCache[userID]; ok {
		for _, val := range D.DataCache[userID] {
			if val.ShortURL == shortURL {
				return nil
			}
		}
	}
	D.DataCache[userID] = append(D.DataCache[userID], modelURL)
	return nil
}

func (D *DB) Close() {
	D.DataCache = nil
}
