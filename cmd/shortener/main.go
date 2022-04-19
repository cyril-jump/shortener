package main

import (
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//db
	url := storage.NewURL()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Routes
	e.GET("/:id", handlers.GetURL(url))
	e.POST("/", handlers.PostURL(url))

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
