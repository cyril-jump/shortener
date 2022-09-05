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

// Flags struct
var Flags struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	ConfigJSON      string
	EnableHTTPS     bool
}

// EnvVar config vars
var EnvVar struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080" json:"server_address,omitempty"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url,omitempty"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	ConfigJSON      string `env:"CONFIG"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https,omitempty"`
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
		json.Unmarshal(appConfigBytes, &EnvVar)
	}

	if srvAddr != "" && EnvVar.ServerAddress == "" {
		cfg["server_address_str"] = srvAddr
	} else {
		cfg["server_address_str"] = EnvVar.ServerAddress
	}

	if hostName != "" && EnvVar.BaseURL == "" {
		cfg["base_url_str"] = hostName
	} else {
		cfg["base_url_str"] = EnvVar.BaseURL
	}

	if fileStoragePath != "" && EnvVar.FileStoragePath == "" {
		cfg["file_storage_path_str"] = fileStoragePath
	} else {
		cfg["file_storage_path_str"] = EnvVar.FileStoragePath
	}

	if databaseDSN != "" && EnvVar.DatabaseDSN == "" {
		cfg["database_dsn_str"] = databaseDSN
	} else {
		cfg["database_dsn_str"] = EnvVar.DatabaseDSN
	}

	if enableHTTPS && !EnvVar.EnableHTTPS {
		cfg["enable_https"] = strconv.FormatBool(enableHTTPS)
	} else {
		cfg["enable_https"] = strconv.FormatBool(EnvVar.EnableHTTPS)
	}

	return &Config{
		cfg: cfg,
	}
}
