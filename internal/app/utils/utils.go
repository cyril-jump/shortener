package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func CheckErr(err error, text string) {
	if err != nil {
		log.Fatal(text, ": ", err)
	}
}

func CreateCookie(c echo.Context, usr storage.Users) string {
	userID := uuid.New().String()
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Value, _ = usr.CreateToken(userID)
	cookie.Name = "cookie"
	c.SetCookie(cookie)
	c.Request().AddCookie(cookie)
	return userID
}
