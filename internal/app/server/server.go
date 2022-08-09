package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo-contrib/pprof"

	"github.com/cyril-jump/shortener/internal/app/handlers"
	"github.com/cyril-jump/shortener/internal/app/middlewares"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

func InitSrv(db storage.DB, cfg storage.Cfg, usr storage.Users, inWorker storage.InWorker) *echo.Echo {

	//server
	srv := handlers.New(db, cfg, usr, inWorker)

	//new Echo instance
	e := echo.New()
	pprof.Register(e)

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
	e.DELETE("/api/user/urls", srv.DelURLsBATCH)

	return e
}
