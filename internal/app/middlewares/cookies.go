package middlewares

import (
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
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
		userID := uuid.New().String()

		_, err := c.Cookie("userID")
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = "userID"
			cookie.Value = userID
			c.SetCookie(cookie)
			c.Request().AddCookie(cookie)
		}

		cookie, err := c.Cookie("cookie")
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = userID
			cookie.Value, _ = M.users.CreateCookie(userID)
			cookie.Path = "/"
			c.SetCookie(cookie)
			c.Request().AddCookie(cookie)

		} else {
			if ok := M.users.CheckCookie(cookie.Value, userID); !ok {
				log.Println(ok)
				cookie := new(http.Cookie)
				cookie.Name = "cookie"
				cookie.Value, _ = M.users.CreateCookie(userID)
				cookie.Path = "/"
				c.SetCookie(cookie)
				c.Request().AddCookie(cookie)
			}
		}

		return next(c)
	}
}
