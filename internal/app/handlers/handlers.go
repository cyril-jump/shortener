package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/middlewares"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"github.com/labstack/echo/v4"
	"io"
	"log"
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
		shortURL, baseURL string
	)
	var userID string

	if id := c.Request().Context().Value(middlewares.UserIDCtxName.String()); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		return c.NoContent(http.StatusNoContent)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url")
	utils.CheckErr(err, "base_url")

	shortURL = utils.Hash(body, hostName)
	baseURL = string(body)

	if err := s.db.SetShortURL(userID, shortURL, baseURL); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.String(http.StatusConflict, shortURL)
		}
		return c.NoContent(http.StatusBadRequest)
	}

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

	var userID string

	if id := c.Request().Context().Value(middlewares.UserIDCtxName.String()); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		return c.NoContent(http.StatusNoContent)
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
	if err = s.db.SetShortURL(userID, response.ShortURL, request.BaseURL); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, response)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, response)
}

func (s Server) GetURLsByUserID(c echo.Context) error {

	var URLs []storage.ModelURL
	var err error

	var userID string

	if id := c.Request().Context().Value(middlewares.UserIDCtxName.String()); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		return c.NoContent(http.StatusNoContent)
	}

	if URLs, err = s.db.GetAllURLsByUserID(userID); err != nil || URLs == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, URLs)
}

func (s Server) PostURLsBATCH(c echo.Context) error {
	var request []storage.ModelURLBatchRequest
	var response []storage.ModelURLBatchResponse
	var model storage.ModelURLBatchResponse

	userID := fmt.Sprintf("%v", c.Request().Context().Value(middlewares.UserIDCtxName.String()))

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url")
	utils.CheckErr(err, "base_url")

	for _, val := range request {
		model.CorID = val.CorID
		model.ShortURL = utils.Hash([]byte(val.BaseURL), hostName)
		response = append(response, model)
		if err := s.db.SetShortURL(userID, model.ShortURL, val.BaseURL); err != nil {
			log.Println("Failed to write to DB: ", err)
		}
	}

	return c.JSON(http.StatusCreated, response)
}

func (s Server) PingDB(c echo.Context) error {

	if err := s.db.Ping(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
