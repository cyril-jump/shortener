package handlers

import (
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_PostURL(t *testing.T) {
	type args struct {
		db        *storage.DB
		cfg       *config.Config
		valueBody string
	}
	tests := []struct {
		name     string
		wantCode int
		args     args
	}{
		{
			name:     "Test PostURL Code 201",
			wantCode: http.StatusCreated,
			args: args{
				db:        storage.NewDB(),
				cfg:       config.NewConfig(":8080", "http://localhost:8080/"),
				valueBody: "https://www.yandex.ru",
			},
		},
		{
			name:     "Test PostURL Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:        storage.NewDB(),
				cfg:       config.NewConfig(":8080", "http://localhost:8080/"),
				valueBody: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := New(tt.args.db, tt.args.cfg)

			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPost, "http://localhost:8080", strings.NewReader(tt.args.valueBody),
			)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/")

			handler := srv.PostURL(c)

			if assert.NoError(t, handler) {
				assert.Equal(t, tt.wantCode, rec.Code)
			}

		})
	}
}

func TestServer_GetURL(t *testing.T) {
	type args struct {
		db       *storage.DB
		cfg      *config.Config
		baseURL  string
		shortURL string
		paramID  string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name:     "Test GetURL Code 307",
			wantCode: http.StatusTemporaryRedirect,
			args: args{
				db:       storage.NewDB(),
				cfg:      config.NewConfig(":8080", "http://localhost:8080"),
				baseURL:  "https://www.yandex.ru",
				shortURL: "http://localhost:8080/f845599b098517893fc2712d32774f53",
				paramID:  "f845599b098517893fc2712d32774f53",
			},
		},
		{
			name:     "Test PostURL Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:       storage.NewDB(),
				cfg:      config.NewConfig(":8080", "http://localhost:8080"),
				baseURL:  "https://www.yandex.ru",
				shortURL: "http://localhost:8080/f845599b098517893fc2712d32774f53",
				paramID:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := New(tt.args.db, tt.args.cfg)
			srv.db.SetURL(tt.args.shortURL, tt.args.baseURL)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.args.paramID)

			handler := srv.GetURL(c)

			if assert.NoError(t, handler) {
				assert.Equal(t, tt.wantCode, rec.Code)
			}

		})
	}
}

func TestServer_PostURLJSON(t *testing.T) {
	type args struct {
		db            *storage.DB
		cfg           *config.Config
		valueBodyJSON string
	}
	tests := []struct {
		name     string
		wantCode int
		args     args
	}{
		{
			name:     "Test PostURLJSON Code 201",
			wantCode: http.StatusCreated,
			args: args{
				db:            storage.NewDB(),
				cfg:           config.NewConfig(":8080", "http://localhost:8080/"),
				valueBodyJSON: `{"url": "https://www.yandex.ru"}`,
			},
		},
		{
			name:     "Test PostURLJSON Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:            storage.NewDB(),
				cfg:           config.NewConfig(":8080", "http://localhost:8080/"),
				valueBodyJSON: `{"url": ""}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := New(tt.args.db, tt.args.cfg)

			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPost, "http://localhost:8080", strings.NewReader(tt.args.valueBodyJSON),
			)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath("/api/shorten")

			handler := srv.PostURLJSON(c)

			if assert.NoError(t, handler) {
				assert.Equal(t, tt.wantCode, rec.Code)
				//assert.Equal(t, userJSON, rec.Body.String())
			}

		})
	}
}
