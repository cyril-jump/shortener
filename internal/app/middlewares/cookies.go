package middlewares

import (
	"context"
	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/labstack/echo/v4"
	"log"
)

type MW struct {
	users storage.Users
}

func New(users storage.Users) *MW {
	return &MW{
		users: users,
	}
}

func (M *MW) SessionWithCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var userID string
		var ok bool

		cookie, err := c.Cookie(config.CookieKey.String())
		if err != nil {
			userID = utils.CreateCookie(c, M.users)
		} else {
			userID, ok = M.users.CheckToken(cookie.Value)
			if !ok {
				userID = utils.CreateCookie(c, M.users)
			}
		}

		log.Println(userID, "  USER1")

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), config.CookieKey.String(), userID)))

		return next(c)
	}
}
