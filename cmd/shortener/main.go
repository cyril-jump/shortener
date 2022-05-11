package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/interfaces"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/rom"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	flag "github.com/spf13/pflag"
	"log"
)

var flags struct {
	a string
	b string
	f string
}

var envVar struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func init() {
	//evn vars
	err := env.Parse(&envVar)
	if err != nil {
		log.Fatal(err)
	}

	//flag
	flag.StringVar(&flags.a, "a", envVar.ServerAddress, "server address")
	flag.StringVar(&flags.b, "b", envVar.BaseURL, "base url")
	flag.StringVar(&flags.f, "f", envVar.FileStoragePath, "file storage path")
	flag.Parse()
}

func main() {

	var err error
	//db
	var db interfaces.Storage

	//config
	cfg := config.NewConfig(flags.a, flags.b, flags.f)

	if cfg.FileStoragePath() != "" {
		db, err = rom.NewDB(cfg.FileStoragePath())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db = ram.NewDB()
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
