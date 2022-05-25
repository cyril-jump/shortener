package config

import (
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/tidwall/gjson"
)

// flags

var Flags struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

// env vars

var EnvVar struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"postgres://dmosk:dmosk@localhost:5432/dmosk?sslmode=disable"`
}

// config

type Config struct {
	serverAddress   string
	baseURL         string
	fileStoragePath string
	databaseDSN     string
}

func (c Config) Get(key string) (string, error) {
	conf := &struct {
		ServerAddress   string `json:"server_address"`
		BaseURL         string `json:"base_url"`
		FileStoragePath string `json:"file_storage_path"`
		DatabaseDSN     string `json:"database_dsn"`
	}{
		ServerAddress:   c.serverAddress,
		BaseURL:         c.baseURL,
		FileStoragePath: c.fileStoragePath,
		DatabaseDSN:     c.databaseDSN,
	}
	buf, err := ffjson.Marshal(conf)
	if err != nil {
		return "", err
	}

	if !gjson.GetBytes(buf, key).Exists() {
		return "", errs.ErrNotFound
	}

	return gjson.GetBytes(buf, key).String(), nil
}

//constructor

func NewConfig(srvAddr, hostName, fileStoragePath, databaseDSN string) *Config {
	return &Config{
		serverAddress:   srvAddr,
		baseURL:         hostName,
		fileStoragePath: fileStoragePath,
		databaseDSN:     databaseDSN,
	}
}

// config interface

type Cfg interface {
	Get(key string) (string, error)
}
