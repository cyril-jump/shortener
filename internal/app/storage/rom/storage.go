package rom

import (
	"bufio"
	"encoding/json"
	"github.com/cyril-jump/shortener/internal/app/interfaces"
	"log"
	"os"
)

// DB

type DB struct {
	file    *os.File
	cache   map[string]string
	encoder *json.Encoder
}

//constructor

func NewDB(filepath string) (*DB, error) {

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)

	if stat, _ := file.Stat(); stat.Size() != 0 {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			err := json.Unmarshal(scanner.Bytes(), &data)
			if err != nil {
				log.Fatal("DB file is damaged.")
			}
		}
	}

	//defer file.Close()

	return &DB{
		file:    file,
		cache:   data,
		encoder: json.NewEncoder(file),
	}, nil
}

func (D *DB) Close() error {
	D.cache = nil
	return D.file.Close()
}

func (D *DB) GetBaseURL(key string) (string, error) {
	if v, ok := D.cache[key]; ok {
		return v, nil
	}
	return "", interfaces.ErrNotFound
}

func (D *DB) SetShortURL(key string, value string) error {
	if _, ok := D.cache[key]; ok {
		return interfaces.ErrAlreadyExists
	}
	D.cache[key] = value
	data := make(map[string]string)
	data[key] = value
	return D.encoder.Encode(&data)
}
