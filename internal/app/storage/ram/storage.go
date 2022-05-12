package ram

import (
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// DB

type DB struct {
	storageURL map[string]string
}

//constructor

func NewDB() *DB {
	return &DB{storageURL: make(map[string]string)}
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	if v, ok := D.storageURL[shortURL]; ok {
		return v, nil
	}
	return "", errs.NotFound
}

func (D *DB) SetShortURL(shortURL, baseURL string) error {
	if _, ok := D.storageURL[shortURL]; ok {
		return errs.AlreadyExists
	}
	D.storageURL[shortURL] = baseURL
	return nil
}

func (D *DB) Close() {
	D.storageURL = nil
}
