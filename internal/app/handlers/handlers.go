package handlers

import (
	"encoding/json"
	"github.com/cyril-jump/shortener/internal/app/interfaces"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Server struct {
	db  interfaces.Storage
	cfg interfaces.Config
}

func New(storage interfaces.Storage, config interfaces.Config) *Server {
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

	shortURL = utils.Hash(body, s.cfg.HostName())
	baseURL = string(body)

	err = s.db.SetShortURL(shortURL, baseURL)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.String(http.StatusCreated, shortURL)
}

func (s Server) GetURL(c echo.Context) error {
	var (
		shortURL, baseURL string
	)

	if c.Param("id") == "" {
		return c.NoContent(http.StatusBadRequest)
	} else {
		shortURL = s.cfg.HostName() + "/" + c.Param("id")
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

	/*	if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	*/
	response.ShortURL = utils.Hash([]byte(request.BaseURL), s.cfg.HostName())
	err = s.db.SetShortURL(response.ShortURL, request.BaseURL)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, response)
}
