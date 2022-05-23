package handlers

import (
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/users"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_PostURL(t *testing.T) {
	type args struct {
		db        *ram.DB
		cfg       *config.Config
		usr       *users.DBUsers
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
				db:        ram.NewDB(),
				cfg:       config.NewConfig(":8080", "http://localhost:8080/", ""),
				usr:       users.New(),
				valueBody: "https://www.yandex.ru",
			},
		},
		{
			name:     "Test PostURL Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:        ram.NewDB(),
				cfg:       config.NewConfig(":8080", "http://localhost:8080/", ""),
				usr:       users.New(),
				valueBody: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := New(tt.args.db, tt.args.cfg, tt.args.usr)

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
		db       *ram.DB
		cfg      *config.Config
		usr      *users.DBUsers
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
				db:       ram.NewDB(),
				cfg:      config.NewConfig(":8080", "http://localhost:8080", ""),
				usr:      users.New(),
				baseURL:  "https://www.yandex.ru",
				shortURL: "http://localhost:8080/f845599b098517893fc2712d32774f53",
				paramID:  "f845599b098517893fc2712d32774f53",
			},
		},
		{
			name:     "Test PostURL Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:       ram.NewDB(),
				cfg:      config.NewConfig(":8080", "http://localhost:8080", ""),
				usr:      users.New(),
				baseURL:  "https://www.yandex.ru",
				shortURL: "http://localhost:8080/f845599b098517893fc2712d32774f53",
				paramID:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			srv := New(tt.args.db, tt.args.cfg, tt.args.usr)

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			tt.args.usr.SetUserID(c.Request().RemoteAddr)
			userID, _ := tt.args.usr.GetUserID(c.Request().RemoteAddr)
			_ = srv.db.SetShortURL(userID, tt.args.shortURL, tt.args.baseURL)
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
		db            *ram.DB
		cfg           *config.Config
		usr           *users.DBUsers
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
				db:            ram.NewDB(),
				cfg:           config.NewConfig(":8080", "http://localhost:8080/", ""),
				usr:           users.New(),
				valueBodyJSON: `{"url": "https://www.yandex.ru"}`,
			},
		},
		{
			name:     "Test PostURLJSON Code 400",
			wantCode: http.StatusBadRequest,
			args: args{
				db:            ram.NewDB(),
				cfg:           config.NewConfig(":8080", "http://localhost:8080/", ""),
				usr:           users.New(),
				valueBodyJSON: `{"url": ""}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srv := New(tt.args.db, tt.args.cfg, tt.args.usr)

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

func TestServer_GetURLsByUserID(t *testing.T) {
	type args struct {
		db        *ram.DB
		cfg       *config.Config
		usr       *users.DBUsers
		isWrite   bool
		baseURL1  string
		shortURL1 string
		baseURL2  string
		shortURL2 string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name:     "Test GetURL Code 200",
			wantCode: http.StatusOK,
			args: args{
				db:        ram.NewDB(),
				cfg:       config.NewConfig(":8080", "http://localhost:8080", ""),
				usr:       users.New(),
				isWrite:   true,
				baseURL1:  "https://www.yandex.ru",
				shortURL1: "http://localhost:8080/f845599b098517893fc2712d32774f53",
				baseURL2:  "https://www.vk.com",
				shortURL2: "http://localhost:9090/15ba5d5d871df48f3b5132ba8c213d23",
			},
		},
		{
			name:     "Test PostURL Code 204",
			wantCode: http.StatusNoContent,
			args: args{
				db:      ram.NewDB(),
				cfg:     config.NewConfig(":8080", "http://localhost:8080", ""),
				usr:     users.New(),
				isWrite: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			srv := New(tt.args.db, tt.args.cfg, tt.args.usr)

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			tt.args.usr.SetUserID("user")
			userID, _ := tt.args.usr.GetUserID("user")
			if tt.args.isWrite {
				_ = srv.db.SetShortURL(userID, tt.args.shortURL1, tt.args.baseURL1)
				_ = srv.db.SetShortURL(userID, tt.args.shortURL2, tt.args.baseURL2)
			}

			handler := srv.GetURLsByUserID(c)

			if assert.NoError(t, handler) {
				assert.Equal(t, tt.wantCode, rec.Code)
			}

		})
	}
}
