package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/middlewares"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/storage/postgres"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/rom"
	"github.com/cyril-jump/shortener/internal/app/storage/users"
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
	flag.StringVarP(&config.Flags.DatabaseDSN, "psqlConn", "d", config.EnvVar.DatabaseDSN, "database URL conn")
	flag.Parse()

}

func main() {

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	//db
	var db storage.DB

	//config
	cfg := config.NewConfig(config.Flags.ServerAddress, config.Flags.BaseURL, config.Flags.FileStoragePath, config.Flags.DatabaseDSN)

	psqlConn, err := cfg.Get("database_dsn")
	utils.CheckErr(err, "")

	fileStoragePath, err := cfg.Get("file_storage_path")
	utils.CheckErr(err, "file_storage_path")

	if fileStoragePath != "" {
		db, err = rom.NewDB(fileStoragePath)
		utils.CheckErr(err, "")
	} else if psqlConn != "" {
		db = postgres.New(psqlConn)
	} else {
		db = ram.NewDB()
	}
	usr := users.New()
	//server
	srv := handlers.New(db, cfg, usr)

	//new Echo instance
	e := echo.New()

	// Middleware
	mw := middlewares.New(usr)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.Use(mw.SessionWithCookies)
	//Routes
	e.GET("/:urlID", srv.GetURL)
	e.GET("/api/user/urls", srv.GetURLsByUserID)
	e.GET("/ping", srv.PingDB)
	e.POST("/", srv.PostURL)
	e.POST("/api/shorten", srv.PostURLJSON)
	e.POST("/api/shorten/batch", srv.PostURLsBATCH)

	// Start Server

	serverAddress, err := cfg.Get("server_address")
	utils.CheckErr(err, "server_address")

	go func() {
		if err = e.Start(serverAddress); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}

	}()

	<-signalChan

	log.Println("Shutting down...")

	cancel()
	if err = e.Shutdown(ctx); err != nil && err != ctx.Err() {
		e.Logger.Fatal(err)
	}

	if err = db.Close(); err != nil {
		log.Fatal(err)
	}
}
