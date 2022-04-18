package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func PostUrl(url *storage.Url) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil || len(body) == 0 {
			return c.NoContent(http.StatusBadRequest)
		}
		shortUrl := hash(body)

		url.Short[shortUrl] = string(body)

		return c.String(http.StatusCreated, shortUrl)
	}
}

func GetUrl(url *storage.Url) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Param("id") == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		shortUrl := "localhost:8080/" + c.Param("id")

		if url.Short[shortUrl] == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set("Location", url.Short[shortUrl])
		return c.NoContent(http.StatusTemporaryRedirect)
	}
}

func hash(url []byte) string {
	host := "localhost:8080/"
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%x", host, hash)
}
