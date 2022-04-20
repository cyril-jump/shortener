package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

// Const and Var

var shortURL = ""

const host = "http://localhost:8080/"

// Handlers

func PostURL(url *storage.URL) echo.HandlerFunc {
	return func(c echo.Context) error {

		body, err := io.ReadAll(c.Request().Body)

		if err != nil || len(body) == 0 {
			return c.NoContent(http.StatusBadRequest)
		} else {
			shortURL = hash(body)
		}

		url.Short[shortURL] = string(body)

		return c.String(http.StatusCreated, shortURL)
	}
}

func GetURL(db *storage.URL) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Param("id") == "" {
			return c.NoContent(http.StatusBadRequest)
		} else {
			shortURL = host + c.Param("id")
		}

		if db.Short[shortURL] == "" {
			return c.NoContent(http.StatusBadRequest)
		} else {
			c.Response().Header().Set("Location", db.Short[shortURL])
			return c.NoContent(http.StatusTemporaryRedirect)
		}

	}
}

// other func
func hash(url []byte) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%x", host, hash)
}
