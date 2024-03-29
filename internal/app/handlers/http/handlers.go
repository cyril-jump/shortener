package http

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// Server struct
type Server struct {
	db       storage.DB
	cfg      storage.Cfg
	usr      storage.Users
	inWorker storage.InWorker
}

// New Server constructor
func New(db storage.DB, config storage.Cfg, usr storage.Users, inWorker storage.InWorker) *Server {
	return &Server{
		db:       db,
		cfg:      config,
		usr:      usr,
		inWorker: inWorker,
	}
}

// Handlers

// PostURL  Accepts a URL string in the request body for shortening
func (s Server) PostURL(c echo.Context) error {

	shortURL := ""
	baseURL := ""
	userID := ""

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

// GetURL Accepts the identifier of the short URL as a URL parameter and returns a response
func (s Server) GetURL(c echo.Context) error {

	shortURL := ""
	baseURL := ""
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

// PostURLJSON Accepting a JSON object in the request body and returning a JSON objec in response
func (s Server) PostURLJSON(c echo.Context) error {
	request := dto.ModelRequestURL{}
	response := dto.ModelResponseURL{}

	userID := ""

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

	response.ShortURL = utils.Hash([]byte(request.BaseURL), hostName)
	if err = s.db.SetShortURL(userID, response.ShortURL, request.BaseURL); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, response)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, response)
}

// GetURLsByUserID Return to the user all ever saved by him
func (s Server) GetURLsByUserID(c echo.Context) error {

	var urls []dto.ModelURL
	var err error
	userID := ""

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}

	urls, err = s.db.GetAllURLsByUserID(userID)
	if err != nil || urls == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, urls)
}

// PostURLsBATCH Accepting in the request body a set of URLs for shortening in the format
func (s Server) PostURLsBATCH(c echo.Context) error {
	request := make([]dto.ModelURLBatchRequest, 0, 20000)
	response := make([]dto.ModelURLBatchResponse, 0, 20000)
	model := dto.ModelURLBatchResponse{}
	userID := ""

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

// DelURLsBATCH Accepts a list of abbreviated URL IDs to delete
func (s Server) DelURLsBATCH(c echo.Context) error {
	userID := ""
	hostName, err := s.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

	if id := c.Request().Context().Value(config.CookieKey); id != nil {
		userID = id.(string)
	}

	if userID == "" {
		userID = utils.CreateCookie(c, s.usr)
	}
	log.Println(userID, "userID")
	model := dto.Task{}
	model.ID = userID

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	deleteURLs := make([]string, 20000)
	err = json.Unmarshal(body, &deleteURLs)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	for _, url := range deleteURLs {
		model.ShortURL = hostName + "/" + url
		s.inWorker.Do(model)
	}

	return c.NoContent(http.StatusAccepted)
}

// PingDB Checks the connection to the database
func (s Server) PingDB(c echo.Context) error {

	if err := s.db.Ping(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

//GetStats get stats
func (s Server) GetStats(c echo.Context) error {
	trustedSubnet, err := s.cfg.Get("trusted_subnet")
	utils.CheckErr(err, "trusted_subnet")

	ip := c.Request().Header.Get("X-Real-IP")
	err = utils.CheckIP(ip, trustedSubnet)
	if err != nil {
		if errors.Is(err, errs.ErrNetNotTrusted) {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	stat, err := s.db.GetStats()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, stat)
}
