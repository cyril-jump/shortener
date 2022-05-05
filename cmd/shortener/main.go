package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {

	//evn var
	envVar := config.EnvVar{}
	err := env.Parse(&envVar)
	if err != nil {
		log.Fatal(err)
	}

	//config
	cfg := config.NewConfig(envVar.ServerAddress, envVar.BaseURL)
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
	e.POST("/api/shorten", srv.PostURLJSON)

	// Start Server
	if err := e.Start(cfg.SrvAddr()); err != nil {
		e.Logger.Fatal(err)
	}
}
