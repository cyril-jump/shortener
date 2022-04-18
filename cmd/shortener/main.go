package main

import (
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
)

const servAdr = "localhost:8080"

func main() {

	url := storage.NewUrl()

	e := echo.New()
	e.GET("/:id", handlers.GetUrl(url))
	e.POST("/", handlers.PostUrl(url))
	e.Logger.Info(e.Start(servAdr))
}
