package handlers

import (
	"encoding/json"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Server struct {
	db  storage.DB
	cfg config.Cfg
	usr storage.Users
}

func New(db storage.DB, config config.Cfg, usr storage.Users) *Server {
	return &Server{
		db:  db,
		cfg: config,
		usr: usr,
	}
}

// Handlers

func (s Server) PostURL(c echo.Context) error {
	var (
		shortURL, baseURL, userID string
	)

	cookie, _ := c.Request().Cookie("cookie")
	userID, _ = s.usr.CheckCookie(cookie.Value)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url")
	utils.CheckErr(err, "base_url")

	shortURL = utils.Hash(body, hostName)
	baseURL = string(body)

	s.db.SetShortURL(userID, shortURL, baseURL)

	return c.String(http.StatusCreated, shortURL)
}

func (s Server) GetURL(c echo.Context) error {
	var (
		shortURL, baseURL string
	)
	var err error

	if c.Param("urlID") == "" {
		return c.NoContent(http.StatusBadRequest)
	} else {
		hostName, err := s.cfg.Get("base_url")
		utils.CheckErr(err, "base_url")
		shortURL = hostName + "/" + c.Param("urlID")
	}

	if baseURL, err = s.db.GetBaseURL(shortURL); err != nil {
		return c.NoContent(http.StatusBadRequest)
	} else {
		c.Response().Header().Set("Location", baseURL)
		return c.NoContent(http.StatusTemporaryRedirect)
	}
}

func (s Server) PostURLJSON(c echo.Context) error {
	var request struct {
		BaseURL string `json:"url"`
	}

	var response struct {
		ShortURL string `json:"result"`
	}

	var (
		userID string
	)

	cookie, _ := c.Request().Cookie("cookie")
	userID, _ = s.usr.CheckCookie(cookie.Value)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if request.BaseURL == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url")
	utils.CheckErr(err, "base_url")

	//userID, _ = s.usr.GetUserID(userName)

	response.ShortURL = utils.Hash([]byte(request.BaseURL), hostName)
	err = s.db.SetShortURL(userID, response.ShortURL, request.BaseURL)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, response)
}

func (s Server) GetURLsByUserID(c echo.Context) error {

	var (
		userID string
	)
	var URLs []storage.ModelURL
	var err error

	cookie, _ := c.Request().Cookie("cookie")
	userID, _ = s.usr.CheckCookie(cookie.Value)

	if URLs, err = s.db.GetAllURLsByUserID(userID); err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, URLs)
}

func (s Server) PingDB(c echo.Context) error {

	if err := s.db.Ping(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
