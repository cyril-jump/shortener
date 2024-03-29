package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// User ID
type User string

func (c User) String() string {
	return string(c)
}

var (
	UserID User = "UserID"
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
	TrustedSubnet   string
}

// EnvVar config vars
type EnvVar struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080" json:"server_address,omitempty"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url,omitempty"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	ConfigJSON      string `env:"CONFIG" envDefault:""`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https,omitempty"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet,omitempty"`
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
func NewConfig(srvAddr, hostName, fileStoragePath, databaseDSN, configJSON string, enableHTTPS bool, trustedSubnet string) *Config {
	cfg := make(map[string]string)
	var appConfig EnvVar

	if configJSON != "" {
		configFile, _ := os.Open(configJSON)
		defer configFile.Close()
		reader := bufio.NewReader(configFile)
		stat, _ := configFile.Stat()
		var appConfigBytes = make([]byte, stat.Size())
		reader.Read(appConfigBytes)
		json.Unmarshal(appConfigBytes, &appConfig)
	}

	if srvAddr != "" && appConfig.ServerAddress == "" {
		cfg["server_address_str"] = srvAddr
	} else {
		cfg["server_address_str"] = appConfig.ServerAddress
	}

	if hostName != "" && appConfig.BaseURL == "" {
		cfg["base_url_str"] = hostName
	} else {
		cfg["base_url_str"] = appConfig.BaseURL
	}

	if fileStoragePath != "" && appConfig.FileStoragePath == "" {
		cfg["file_storage_path_str"] = fileStoragePath
	} else {
		cfg["file_storage_path_str"] = appConfig.FileStoragePath
	}

	if databaseDSN != "" && appConfig.DatabaseDSN == "" {
		cfg["database_dsn_str"] = databaseDSN
	} else {
		cfg["database_dsn_str"] = appConfig.DatabaseDSN
	}

	if enableHTTPS && !appConfig.EnableHTTPS {
		cfg["enable_https"] = strconv.FormatBool(enableHTTPS)
	} else {
		cfg["enable_https"] = strconv.FormatBool(appConfig.EnableHTTPS)
	}

	if trustedSubnet != "" && appConfig.TrustedSubnet == "" {
		cfg["trusted_subnet"] = trustedSubnet
	} else {
		cfg["trusted_subnet"] = trustedSubnet
	}

	return &Config{
		cfg: cfg,
	}
}
