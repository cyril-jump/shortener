package middlewares

import (
	"github.com/cyril-jump/shortener/internal/app/storage"
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
		M.users.SetUserID("user")
		userID, _ := M.users.GetUserID("user")

		cookie, err := c.Cookie("cookie")
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = "cookie"
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
