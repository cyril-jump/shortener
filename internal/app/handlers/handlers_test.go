package handlers

import (
	"context"
	"encoding/json"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/middlewares"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/users"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type (
	// for_TestServer_PostURLJSON
	RequestURL struct {
		URL string `json:"url"`
	}
)

type Suite struct {
	suite.Suite
	db      *ram.DB
	cfg     *config.Config
	usr     *users.DBUsers
	e       *echo.Echo
	router  *echo.Router
	testSrv *httptest.Server
	mw      *middlewares.MW
	srv     *Server
	ctx     context.Context
}

func (suite *Suite) SetupTest() {
	suite.e = echo.New()
	suite.router = echo.NewRouter(suite.e)
	suite.db = ram.NewDB(suite.ctx)
	suite.cfg = config.NewConfig(
		":8080",
		"http://localhost:8080",
		"",
		"postgres://dmosk:dmosk@localhost:5432/dmosk?sslmode=disable")
	suite.usr = users.New(suite.ctx)
	suite.router = echo.NewRouter(suite.e)
	suite.testSrv = httptest.NewServer(suite.e)
	suite.mw = middlewares.New(suite.usr)
	suite.srv = New(suite.db, suite.cfg, suite.usr)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) TestServer_PostURL() {

	suite.e.Use(suite.mw.SessionWithCookies)
	suite.e.POST("/", suite.srv.PostURL)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "Test PostURL Code 201",
			URL:  "https://www.yandex.ru",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Test PostURL Code 400",
			URL:  "",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.URL)
			client := resty.New()
			res, err := client.R().SetBody(payload).Post(suite.testSrv.URL)
			if err != nil {
				t.Fatalf("Could not create POST request")
			}
			assert.Equal(t, tt.want.code, res.StatusCode())
		})
	}
	defer suite.testSrv.Close()
}

func (suite *Suite) TestServer_PostURLJSON() {

	suite.e.Use(suite.mw.SessionWithCookies)
	suite.e.POST("/api/shorten", suite.srv.PostURLJSON)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		URL  RequestURL
		want want
	}{
		{
			name: "Test PostURLJSON Code 201",
			URL: RequestURL{
				URL: "https://www.yandex.ru",
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Test PostURLJSON Code 400",
			URL: RequestURL{
				URL: "",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.URL)
			payload := strings.NewReader(string(reqBody))
			client := resty.New()
			res, err := client.R().SetBody(payload).Post(suite.testSrv.URL + "/api/shorten")
			if err != nil {
				t.Fatalf("Could not perform JSON POST request")
			}
			//t.Logf(string(res.Body()))
			assert.Equal(t, tt.want.code, res.StatusCode())
		})
	}
	defer suite.testSrv.Close()

}

func (suite *Suite) TestServer_GetURL() {

	ShortURL := "http://localhost:8080/f845599b098517893fc2712d32774f53"
	BaseURL := "https://www.yandex.ru"
	userID := uuid.New().String()

	suite.e.GET("/:urlID", suite.srv.GetURL)
	_ = suite.db.SetShortURL(userID, ShortURL, BaseURL)

	type want struct {
		code int
	}
	tests := []struct {
		name     string
		ShortURL string
		want     want
	}{
		{
			name:     "Test GetURL Code 307",
			ShortURL: "f845599b098517893fc2712d32774f53",
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
		{
			name:     "Test PostURL Code 400",
			ShortURL: "620f2a73709959c2a511d9be58e2f9ff",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			client := resty.New()
			client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}))

			res, err := client.R().SetPathParams(map[string]string{"urlID": tt.ShortURL}).Get(suite.testSrv.URL + "/{urlID}")
			if err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(t, tt.want.code, res.StatusCode())
		})
	}
	defer suite.testSrv.Close()
}

func (suite *Suite) TestServer_GetURLsByUserID() {
	ShortURL1 := "http://localhost:8080/f845599b098517893fc2712d32774f53"
	BaseURL1 := "https://www.yandex.ru"
	userID1 := uuid.New().String()
	token1, _ := suite.usr.CreateCookie(userID1)
	userID2 := uuid.New().String()
	token2, _ := suite.usr.CreateCookie(userID2)
	_ = suite.db.SetShortURL(userID1, ShortURL1, BaseURL1)
	suite.e.Use(suite.mw.SessionWithCookies)
	suite.e.GET("/api/user/urls", suite.srv.GetURLsByUserID)

	type want struct {
		code int
	}
	tests := []struct {
		name  string
		token string
		want  want
	}{
		{
			name:  "Test GetURLsByUserID Code 200",
			token: token1,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:  "Test GetURLsByUserID Code 204",
			token: token2,
			want: want{
				code: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: tt.token,
				Path:  "/",
			})
			res, err := client.R().Get(suite.testSrv.URL + "/api/user/urls")
			if err != nil {
				t.Fatalf("Could not perform GET by userID request")
			}
			assert.Equal(t, tt.want.code, res.StatusCode())
		})
	}
	defer suite.testSrv.Close()
}
