package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/interfaces"
	"github.com/cyril-jump/shortener/internal/app/storage/storage_ram"
	"github.com/cyril-jump/shortener/internal/app/storage/storage_rom"
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

	//db
	var db interfaces.Storage
	//config
	cfg := config.NewConfig(envVar.ServerAddress, envVar.BaseURL, envVar.FileStoragePath)

	if cfg.FileStoragePath() != "" {
		db, err = storage_rom.NewDB(cfg.FileStoragePath())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db = storage_ram.NewDB()
	}
	defer db.Close()

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
