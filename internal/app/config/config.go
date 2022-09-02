package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// contextKey const
type contextKey string

const (
	CookieKey = contextKey("cookie")
)

func (c contextKey) String() string {
	return string(c)
}

// ConfigJSON struct
type ConfigJSON struct {
	ServerAddress   string `json:"server_address,omitempty"`
	BaseURL         string `json:"base_url,omitempty"`
	FileStoragePath string `json:"file_storage_path,omitempty"`
	DatabaseDSN     string `json:"database_dsn,omitempty"`
	EnableHTTPS     bool   `json:"enable_https,omitempty"`
}

// Flags struct
var Flags struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	ConfigJSON      string
	EnableHTTPS     bool
}

// EnvVar environment vars
var EnvVar struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	ConfigJSON      string `env:"CONFIG"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
}

// Config struct
type Config struct {
	cfg map[string]string
}

// Get Config element by key
func (c Config) Get(key string) (string, error) {
	if _, ok := c.cfg[key]; !ok {
		return "", errs.ErrNotFound
	}
	return c.cfg[key], nil
}

// NewConfig config constructor
func NewConfig(srvAddr, hostName, fileStoragePath, databaseDSN, configJSON string, enableHTTPS bool) *Config {
	cfg := make(map[string]string)

	if configJSON != "" {
		configFile, _ := os.Open(configJSON)
		defer configFile.Close()
		reader := bufio.NewReader(configFile)
		stat, _ := configFile.Stat()
		var appConfigBytes = make([]byte, stat.Size())
		reader.Read(appConfigBytes)
		var appConfig ConfigJSON
		json.Unmarshal(appConfigBytes, &appConfig)
		cfg["server_address_str"] = appConfig.ServerAddress
		cfg["base_url_str"] = appConfig.BaseURL
		cfg["file_storage_path_str"] = appConfig.FileStoragePath
		cfg["database_dsn_str"] = appConfig.DatabaseDSN
		cfg["enable_https"] = strconv.FormatBool(appConfig.EnableHTTPS)
	}

	cfg["server_address_str"] = srvAddr
	cfg["base_url_str"] = hostName
	cfg["file_storage_path_str"] = fileStoragePath
	cfg["database_dsn_str"] = databaseDSN
	cfg["enable_https"] = strconv.FormatBool(enableHTTPS)
	return &Config{
		cfg: cfg,
	}
}
