package utils

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

// Hash Get short URL from base URL
func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
}

// CheckErr Check errors
func CheckErr(err error, text string) {
	if err != nil {
		log.Fatal(text, ": ", err)
	}
}

// CreateCookie Create cookie for user
func CreateCookie(c echo.Context, usr storage.Users) string {
	userID := uuid.New().String()
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Value, _ = usr.CreateToken(userID)
	cookie.Name = config.CookieKey.String()
	c.SetCookie(cookie)
	c.Request().AddCookie(cookie)
	return userID
}

// CheckIP util
func CheckIP(ip string, trustedNet string) error {
	if ip == "" {
		return errs.ErrNetNotTrusted
	}
	ipRequest, _, err := net.ParseCIDR(ip)
	if err != nil {
		return err
	}
	var ipnet *net.IPNet = nil
	if trustedNet != "" {
		_, ipnet, err = net.ParseCIDR(trustedNet)
		if err != nil {
			return err
		}
	}
	if ipnet.Contains(ipRequest) {
		return nil
	}
	return errs.ErrNetNotTrusted
}
