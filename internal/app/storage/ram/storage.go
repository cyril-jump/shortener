package ram

import (
	"context"
	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"log"
)

// DB

type DB struct {
	DataCache map[string][]dto.ModelURL
	GlobalBD  map[string]string
	ctx       context.Context
}

//constructor

func NewDB(ctx context.Context) *DB {

	return &DB{
		DataCache: make(map[string][]dto.ModelURL),
		GlobalBD:  make(map[string]string),
		ctx:       ctx,
	}
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	if v, ok := D.GlobalBD[shortURL]; ok {
		return v, nil
	}
	return "", errs.ErrNoContent
}

func (D *DB) GetAllURLsByUserID(userID string) ([]dto.ModelURL, error) {
	if _, ok := D.DataCache[userID]; ok {
		return D.DataCache[userID], nil
	}
	return nil, errs.ErrNoContent
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	modelURL := dto.ModelURL{
		ShortURL: shortURL,
		BaseURL:  baseURL,
	}

	if _, ok := D.GlobalBD[shortURL]; !ok {
		D.GlobalBD[shortURL] = baseURL
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

func (D *DB) DelBatchShortURLs(t []dto.Task) error {
	log.Println(t)
	return nil
}

func (D *DB) Ping() error {
	return nil
}

func (D *DB) Close() error {
	D.DataCache = nil
	return nil
}
