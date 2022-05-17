package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/rom"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	//evn vars
	err := env.Parse(&config.EnvVar)
	if err != nil {
		log.Fatal(err)
	}

	//flags
	flag.StringVarP(&config.Flags.ServerAddress, "address", "a", config.EnvVar.ServerAddress, "server address")
	flag.StringVarP(&config.Flags.BaseURL, "base", "b", config.EnvVar.BaseURL, "base url")
	flag.StringVarP(&config.Flags.FileStoragePath, "file", "f", config.EnvVar.FileStoragePath, "file storage path")
	flag.Parse()

}

func main() {

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//db
	var db storage.DB

	//config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.BaseURL, config.Flags.FileStoragePath)

	fileStoragePath, err := cfg.Get("file_storage_path")
	utils.CheckErr(err, "file_storage_path")

	if fileStoragePath != "" {
		db, err = rom.NewDB(fileStoragePath)
		utils.CheckErr(err, "")
	} else {
		db = ram.NewDB()
	}

	//server
	srv := handlers.New(db, cfg)

	//new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	//Routes
	e.GET("/:id", srv.GetURL)
	e.POST("/", srv.PostURL)
	e.POST("/api/shorten", srv.PostURLJSON)

	// Start Server

	serverAddress, err := cfg.Get("server_address")
	utils.CheckErr(err, "server_address")

	go func() {
		if err = e.Start(serverAddress); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}

	}()

	for {
		select {
		case <-signalChan:
			log.Println("Shutting down...")
			cancel()
			if err = e.Shutdown(ctx); err != nil && err != ctx.Err() {
				e.Logger.Fatal(err)
			}
			db.Close()
			os.Exit(0)
		}
	}
}
