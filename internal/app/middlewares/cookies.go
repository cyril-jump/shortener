package middlewares

import (
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
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

		cookie, err := c.Cookie("cookie")
		if err != nil {
			userID := uuid.New().String()
			cookie := new(http.Cookie)
			cookie.Path = "/"
			cookie.Value, _ = M.users.CreateCookie(userID)
			cookie.Name = "cookie"
			c.SetCookie(cookie)
			c.Request().AddCookie(cookie)
		} else {
			if _, ok := M.users.CheckCookie(cookie.Value); !ok {
				userID := uuid.New().String()
				cookie := new(http.Cookie)
				cookie.Path = "/"
				cookie.Value, _ = M.users.CreateCookie(userID)
				cookie.Name = "cookie"
				c.SetCookie(cookie)
				c.Request().AddCookie(cookie)
			}
		}

		return next(c)
	}
}
