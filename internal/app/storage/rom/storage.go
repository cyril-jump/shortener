package rom

import (
	"bufio"
	"encoding/json"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"log"
	"os"
)

// DB

type DB struct {
	file      *os.File
	DataFile  ModelFile `json:"data_file"`
	DataCache map[string][]storage.ModelURL
	GlobalBD  map[string]string
	encoder   *json.Encoder
}

type ModelFile struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	BaseURL  string `json:"base_url"`
}

//constructor

func NewDB(filepath string) (*DB, error) {

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	var modelURL storage.ModelURL
	dataCache := make(map[string][]storage.ModelURL)
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
	}, nil
}

func (D *DB) Close() {
	D.DataCache = nil
	if err := D.file.Close(); err != nil {
		log.Fatal(err)
	}
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	if v, ok := D.GlobalBD[shortURL]; ok {
		return v, nil
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

	D.DataFile.UserID = userID
	D.DataFile.ShortURL = shortURL
	D.DataFile.BaseURL = baseURL

	modelURL := storage.ModelURL{
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
