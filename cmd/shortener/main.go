package main

import (
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	//config
	cfg := config.NewConfig(":8080", "http://localhost:8080/")
	//db
	db := storage.NewDB()

	//server
	srv := handlers.New(db, cfg)

	//new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Routes
	e.GET("/:id", srv.GetURL)
	e.POST("/", srv.PostURL)

	// Start Server
	if err := e.Start(cfg.SrvAddr()); err != nil {
		e.Logger.Fatal(err)
	}
}
