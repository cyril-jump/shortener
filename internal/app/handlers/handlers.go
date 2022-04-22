package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

// Handlers

func PostURL(db *storage.DB, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			shortURL, baseURL string
		)

		body, err := io.ReadAll(c.Request().Body)

		if err != nil || len(body) == 0 {
			return c.NoContent(http.StatusBadRequest)
		} else {
			shortURL = hash(body, cfg.HostName())
			baseURL = string(body)
		}

		db.SetURL(shortURL, baseURL)

		return c.String(http.StatusCreated, shortURL)
	}
}

func GetURL(db *storage.DB, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			shortURL, baseURL string
		)

		if c.Param("id") == "" {
			return c.NoContent(http.StatusBadRequest)
		} else {
			shortURL = cfg.HostName() + c.Param("id")
		}

		if baseURL = db.BaseURL(shortURL); baseURL == "" {
			return c.NoContent(http.StatusBadRequest)
		} else {
			c.Response().Header().Set("Location", baseURL)
			return c.NoContent(http.StatusTemporaryRedirect)
		}

	}
}

// other func
func hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%x", hostName, hash)
}
