package handlers

import (
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetURL(t *testing.T) {
	type args struct {
		db  *storage.DB
		cfg *storage.Config
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		params   string
	}{
		{
			name: "Test GetURL Code 307",
			args: args{
				db: &storage.DB{StorageURL: map[string]string{
					"http://localhost:8080/f845599b098517893fc2712d32774f53": "https://www.yandex.ru"}},
				cfg: &storage.Config{
					SrvAddr:  ":8080",
					HostName: "http://localhost:8080/",
				},
			},
			wantCode: http.StatusTemporaryRedirect,
			params:   "f845599b098517893fc2712d32774f53",
		},
		{
			name: "Test PostURL Code 400",
			args: args{db: &storage.DB{StorageURL: map[string]string{
				"http://localhost:8080/f845599b098517893fc2712d32774f53": "https://www.yandex.ru"}},
				cfg: &storage.Config{
					SrvAddr:  ":8080",
					HostName: "http://localhost:8080/",
				},
			},
			wantCode: http.StatusBadRequest,
			params:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.params)
			db := tt.args.db
			cfg := tt.args.cfg
			handler := GetURL(db, cfg)
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, tt.wantCode, rec.Code)
			}
		})
	}
}

func TestPostURL(t *testing.T) {
	type args struct {
		db  *storage.DB
		cfg *storage.Config
	}
	tests := []struct {
		name      string
		valueBody string
		wantCode  int
		args      args
	}{
		{
			name:      "Test PostURL Code 201",
			valueBody: "https://www.yandex.ru",
			wantCode:  http.StatusCreated,
			args: args{
				db: &storage.DB{StorageURL: map[string]string{}},
				cfg: &storage.Config{
					SrvAddr:  ":8080",
					HostName: "http://localhost:8080/",
				},
			},
		},
		{
			name:      "Test PostURL Code 400",
			valueBody: "",
			wantCode:  http.StatusBadRequest,
			args: args{
				db: &storage.DB{StorageURL: map[string]string{}},
				cfg: &storage.Config{
					SrvAddr:  ":8080",
					HostName: "http://localhost:8080/",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", strings.NewReader(tt.valueBody))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/")
			db := tt.args.db
			cfg := tt.args.cfg
			handler := PostURL(db, cfg)
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, tt.wantCode, rec.Code)
			}
		})
	}
}
