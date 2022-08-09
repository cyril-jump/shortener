package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/cyril-jump/shortener/internal/app/dto"
)

const PostURL = "http://localhost:8080/"
const PostURLJSON = "http://localhost:8080/api/shorten"
const GetURL = "http://localhost:8080/"
const GetURLsByUserID = "http://localhost:8080/api/user/urls"
const PostURLsBATCH = "http://localhost:8080/api/shorten/batch"
const DelURLsBATCH = "http://localhost:8080/api/user/urls"

func Benchmark_PostURLsBATCH(b *testing.B) {

	URLs := make([]dto.ModelURLBatchRequest, 20000)
	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "cookie",
		Value: token,
		Path:  "/",
	})

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			URLs = nil
			for j := 0; j < 5000; j++ {
				id := uuid.New().String()
				URLs = append(URLs, dto.ModelURLBatchRequest{BaseURL: "https://www." + id + ".com", CorID: "user"})
			}
			reqBody, _ := json.Marshal(URLs)
			payload := strings.NewReader(string(reqBody))

			b.StartTimer()
			_, _ = client.R().SetBody(payload).Post(PostURLsBATCH)
		}
	})
}

func Benchmark_PostURL(b *testing.B) {

	var URL string
	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)
	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "cookie",
		Value: token,
		Path:  "/",
	})

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			URL = "https://www." + id + ".com"

			reqBody, _ := json.Marshal(URL)
			payload := strings.NewReader(string(reqBody))

			b.StartTimer()
			_, _ = client.R().SetBody(payload).Post(PostURL)
		}
	})
}

func Benchmark_PostURLJSON(b *testing.B) {

	URL := dto.ModelRequestURL{}
	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)
	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "cookie",
		Value: token,
		Path:  "/",
	})

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			URL = dto.ModelRequestURL{BaseURL: "https://www." + id + ".com"}

			reqBody, _ := json.Marshal(URL)
			payload := strings.NewReader(string(reqBody))

			b.StartTimer()
			_, _ = client.R().SetBody(payload).Post(PostURLJSON)
		}
	})
}

func Benchmark_GetURL(b *testing.B) {

	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)
	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "cookie",
		Value: token,
		Path:  "/",
	})

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()

			b.StartTimer()
			_, _ = client.R().Get(GetURLsByUserID)
		}
	})
}

func Benchmark_GetURLsByUserID(b *testing.B) {

	client := resty.New()
	client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}))

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			ShortURL := uuid.New().String()

			b.StartTimer()
			_, _ = client.R().Get(GetURL + ShortURL)
		}
	})
}

func Benchmark_DelURLsBATCH(b *testing.B) {

	ShortURLs := make([]string, 20000)
	userID := "eae635dd-3ad1-47d3-8f6d-f12b0268eea7"
	token, _ := CreateToken(userID)

	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "cookie",
		Value: token,
		Path:  "/",
	})

	b.ResetTimer()
	b.Run("b", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			ShortURLs = nil
			for j := 0; j < 5000; j++ {
				id := uuid.New().String()
				ShortURLs = append(ShortURLs, id)
			}
			reqBody, _ := json.Marshal(ShortURLs)
			payload := strings.NewReader(string(reqBody))

			b.StartTimer()
			_, _ = client.R().SetBody(payload).Delete(DelURLsBATCH)
		}
	})
}

func CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": userID})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString, nil
}
