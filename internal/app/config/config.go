package config

import "github.com/cyril-jump/shortener/internal/app/utils/errs"

// context const
type contextKey string

const (
	CookieKey = contextKey("cookie")
)

func (c contextKey) String() string {
	return string(c)
}

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
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// config

type Config struct {
	cfg map[string]string
}

func (c Config) Get(key string) (string, error) {
	if _, ok := c.cfg[key]; !ok {
		return "", errs.ErrNotFound
	}
	return c.cfg[key], nil
}

//constructor

func NewConfig(srvAddr, hostName, fileStoragePath, databaseDSN string) *Config {
	cfg := make(map[string]string)
	cfg["server_address_str"] = srvAddr
	cfg["base_url_str"] = hostName
	cfg["file_storage_path_str"] = fileStoragePath
	cfg["database_dsn_str"] = databaseDSN
	return &Config{
		cfg: cfg,
	}
}
