package rom

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// DB Rom struct
type DB struct {
	file      *os.File
	DataFile  ModelFile `json:"data_file"`
	DataCache map[string][]dto.ModelURL
	GlobalBD  map[string]string
	encoder   *json.Encoder
	ctx       context.Context
}

// ModelFile File model
type ModelFile struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"base_url"`
}

// NewDB  Rom constructor
func NewDB(ctx context.Context, filepath string) (*DB, error) {

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	var modelURL dto.ModelURL
	dataCache := make(map[string][]dto.ModelURL)
	globalDB := make(map[string]string)
	var dataFile ModelFile

	if stat, _ := file.Stat(); stat.Size() != 0 {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			err := json.Unmarshal(scanner.Bytes(), &dataFile)
			if err != nil {
				log.Fatal("DB file is damaged.", err)
			}
			globalDB[dataFile.ShortURL] = dataFile.BaseURL
			modelURL.ShortURL = dataFile.ShortURL
			modelURL.BaseURL = dataFile.BaseURL
			dataCache[dataFile.UserID] = append(dataCache[dataFile.UserID], modelURL)

		}
	}

	return &DB{
		file:      file,
		DataFile:  dataFile,
		DataCache: dataCache,
		GlobalBD:  globalDB,
		encoder:   json.NewEncoder(file),
		ctx:       ctx,
	}, nil
}

// GetBaseURL Get base URL from file
func (D *DB) GetBaseURL(shortURL string) (string, error) {
	if v, ok := D.GlobalBD[shortURL]; ok {
		return v, nil
	}
	return "", errs.ErrNoContent
}

// GetAllURLsByUserID Get all URLs by UserID from file
func (D *DB) GetAllURLsByUserID(userID string) ([]dto.ModelURL, error) {
	if _, ok := D.DataCache[userID]; ok {
		return D.DataCache[userID], nil
	}
	return nil, errs.ErrNoContent
}

// SetShortURL Set short URL in file
func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {

	D.DataFile.UserID = userID
	D.DataFile.ShortURL = shortURL
	D.DataFile.BaseURL = baseURL

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

	return D.encoder.Encode(&D.DataFile)
}

func (D *DB) DelBatchShortURLs(task []dto.Task) error {
	log.Println(task)
	return nil
}

func (D *DB) Ping() error {
	return nil
}

func (D *DB) Close() error {
	D.DataCache = nil
	return D.file.Close()
}

func (D *DB) GetStats() (dto.Stat, error) {

	return dto.Stat{}, nil
}
