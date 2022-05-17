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
}

func New(storage storage.DB, config config.Cfg) *Server {
	return &Server{
		db:  storage,
		cfg: config,
	}
}

// Handlers

func (s Server) PostURL(c echo.Context) error {
	var (
		shortURL, baseURL string
	)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url")
	utils.CheckErr(err, "base_url")

	shortURL = utils.Hash(body, hostName)
	baseURL = string(body)

	s.db.SetShortURL(shortURL, baseURL)

	return c.String(http.StatusCreated, shortURL)
}

func (s Server) GetURL(c echo.Context) error {
	var (
		shortURL, baseURL string
	)

	if c.Param("id") == "" {
		return c.NoContent(http.StatusBadRequest)
	} else {
		hostName, err := s.cfg.Get("base_url")
		utils.CheckErr(err, "base_url")
		shortURL = hostName + "/" + c.Param("id")
	}

	if baseURL, _ = s.db.GetBaseURL(shortURL); baseURL == "" {
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

	response.ShortURL = utils.Hash([]byte(request.BaseURL), hostName)
	err = s.db.SetShortURL(response.ShortURL, request.BaseURL)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, response)
}
