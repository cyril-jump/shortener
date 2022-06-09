package handlers

import (
	"encoding/json"
	"errors"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
)

type Server struct {
	db       storage.DB
	cfg      storage.Cfg
	usr      storage.Users
	inWorker storage.InWorker
}

func New(db storage.DB, config storage.Cfg, usr storage.Users, inWorker storage.InWorker) *Server {
	return &Server{
		db:       db,
		cfg:      config,
		usr:      usr,
		inWorker: inWorker,
	}
}

// Handlers

func (s Server) PostURL(c echo.Context) error {
	var (
		shortURL, baseURL string
	)
	var userID string

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}

	body, err := io.ReadAll(c.Request().Body)
	log.Println(string(body))
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

	shortURL = utils.Hash(body, hostName)
	baseURL = string(body)

	if err := s.db.SetShortURL(userID, shortURL, baseURL); err != nil {
		log.Println(err)
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
		hostName, err := s.cfg.Get("base_url_str")
		utils.CheckErr(err, "base_url_str")
		shortURL = hostName + "/" + c.Param("urlID")
	}

	if baseURL, err = s.db.GetBaseURL(shortURL); err != nil {
		if errors.Is(err, errs.ErrWasDeleted) {
			return c.NoContent(http.StatusGone)
		} else {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	c.Response().Header().Set("Location", baseURL)
	return c.NoContent(http.StatusTemporaryRedirect)
}

func (s Server) PostURLJSON(c echo.Context) error {
	var request dto.ModelRequestURL
	var response dto.ModelResponseURL

	var userID string

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
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

	hostName, err := s.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

	//userID, _ = s.usr.GetUserID(userName)

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

	var URLs []dto.ModelURL
	var err error

	var userID string

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}

	if URLs, err = s.db.GetAllURLsByUserID(userID); err != nil || URLs == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, URLs)
}

func (s Server) PostURLsBATCH(c echo.Context) error {
	var request []dto.ModelURLBatchRequest
	var response []dto.ModelURLBatchResponse
	var model dto.ModelURLBatchResponse

	var userID string

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	hostName, err := s.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

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

func (s Server) DelURLsBATCH(c echo.Context) error {

	//s.inWorker.Do(model)

	var userID string

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}
	var model dto.Task
	model.Id = userID

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		log.Println(err)
		log.Println(string(body))
		return c.NoContent(http.StatusBadRequest)
	}
	deleteURLs := make([]string, 0)
	err = json.Unmarshal(body, &deleteURLs)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	for _, url := range deleteURLs {
		model.ShortURL = url
		s.inWorker.Do(model)
	}

	return c.NoContent(http.StatusAccepted)
}

func (s Server) PingDB(c echo.Context) error {

	if err := s.db.Ping(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
