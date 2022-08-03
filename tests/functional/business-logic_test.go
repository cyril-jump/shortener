package functionaltest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/storage/users"
)

func Test_PostURL(t *testing.T) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	tests := []struct {
		name     string
		token    string
		addrSrv  string
		URLs     []string
		wantCode int
	}{
		{
			name:    "success",
			addrSrv: "http://localhost:8080/",
			URLs: []string{
				"https://www.yandex.ru",
				"https://www.google.com",
				"https://www.vk.com",
			},
			token:    token,
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: tt.token,
				Path:  "/",
			})

			for _, URL := range tt.URLs {

				payload := strings.NewReader(URL)
				res, err := client.R().SetBody(payload).Post(tt.addrSrv)
				if err != nil {
					t.Fatalf(err.Error())
				}
				assert.Equal(t, tt.wantCode, res.StatusCode())

			}
		})
	}
}

func Test_PostURLJSON(t *testing.T) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	tests := []struct {
		name     string
		addrSrv  string
		token    string
		URLs     []dto.ModelRequestURL
		wantCode int
	}{
		{
			name:    "success",
			addrSrv: "http://localhost:8080/api/shorten",
			token:   token,
			URLs: []dto.ModelRequestURL{
				{BaseURL: "https://www.redis.io"},
				{BaseURL: "https://www.mozila.com"},
				{BaseURL: "https://www.echo.labstack.com"},
			},
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: token,
				Path:  "/",
			})

			for _, URL := range tt.URLs {
				reqBody, _ := json.Marshal(URL)
				payload := strings.NewReader(string(reqBody))

				res, err := client.R().SetBody(payload).Post(tt.addrSrv)
				if err != nil {
					t.Fatalf(err.Error())
				}
				assert.Equal(t, tt.wantCode, res.StatusCode())

			}
		})
	}
}

func Test_GetURL(t *testing.T) {
	tests := []struct {
		name     string
		addrSrv  []string
		wantCode int
	}{
		{
			name: "success",
			addrSrv: []string{
				"http://localhost:8080/f845599b098517893fc2712d32774f53",
				"http://localhost:8080/8ffdefbdec956b595d257f0aaeefd623",
				"http://localhost:8080/15ba5d5d871df48f3b5132ba8c213d23",
				"http://localhost:8080/42d22622bce626bc0ffc5aee090da6be",
				"http://localhost:8080/6c1f6e66a20a4cd2c61643bc0e25b388",
				"http://localhost:8080/1c0bd3ded146caf4d92df952533995f9",
			},
			wantCode: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}))

			for _, addr := range tt.addrSrv {
				res, err := client.R().Get(addr)
				if err != nil {
					t.Fatalf(err.Error())
				}
				assert.Equal(t, tt.wantCode, res.StatusCode())

			}
		})
	}
}

func Test_GetURLsByUserID(t *testing.T) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	tests := []struct {
		name     string
		addrSrv  string
		token    string
		wantCode int
	}{
		{
			name:     "success",
			addrSrv:  "http://localhost:8080/api/user/urls",
			token:    token,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usr := users.New(context.TODO())
			userID1 := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
			token, err := usr.CreateToken(userID1)
			if err != nil {
				t.Fatalf(err.Error())
			}

			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: token,
				Path:  "/",
			})

			res, err := client.R().Get(tt.addrSrv)
			if err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(t, tt.wantCode, res.StatusCode())

		})
	}
}

func Test_PostURLsBATCH(t *testing.T) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	tests := []struct {
		name     string
		addrSrv  string
		token    string
		URLs     []dto.ModelURLBatchRequest
		wantCode int
	}{
		{
			name:    "success",
			addrSrv: "http://localhost:8080/api/shorten/batch",
			token:   token,
			URLs: []dto.ModelURLBatchRequest{
				{BaseURL: "https://www.atlanta.io", CorID: userID},
				{BaseURL: "https://www.fifefox.com", CorID: userID},
				{BaseURL: "https://www.gameover.com", CorID: userID},
			},
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usr := users.New(context.TODO())
			userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
			token, err := usr.CreateToken(userID)
			if err != nil {
				t.Fatalf(err.Error())
			}

			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: token,
				Path:  "/",
			})

			reqBody, _ := json.Marshal(tt.URLs)
			payload := strings.NewReader(string(reqBody))

			res, err := client.R().SetBody(payload).Post(tt.addrSrv)
			if err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(t, tt.wantCode, res.StatusCode())

		})
	}
}

func Test_DelURLsBATCH(t *testing.T) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	tests := []struct {
		name     string
		addrSrv  string
		token    string
		URLs     []string
		wantCode int
	}{
		{
			name:    "success",
			addrSrv: "http://localhost:8080/api/user/urls",
			token:   token,
			URLs: []string{
				"f845599b098517893fc2712d32774f53",
				"8ffdefbdec956b595d257f0aaeefd623",
				"15ba5d5d871df48f3b5132ba8c213d23",
				"42d22622bce626bc0ffc5aee090da6be",
				"6c1f6e66a20a4cd2c61643bc0e25b388",
				"1c0bd3ded146caf4d92df952533995f9",
			},
			wantCode: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "cookie",
				Value: token,
				Path:  "/",
			})

			reqBody, _ := json.Marshal(tt.URLs)
			payload := strings.NewReader(string(reqBody))

			log.Println(string(reqBody))

			res, err := client.R().SetBody(payload).Delete(tt.addrSrv)
			if err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(t, tt.wantCode, res.StatusCode())

		})
	}
}

func CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": userID})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString, nil
}
