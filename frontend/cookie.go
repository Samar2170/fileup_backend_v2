package frontend

import (
	"fileupbackendv2/internal/storage"
	"fileupbackendv2/internal/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var cookieExemptPaths = []string{
	"/app/form/api-key/",
	"/app/verify/api-key/",
}

func ResetCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "api-key"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		cookie, err := c.Cookie("api-key")
		if err != nil || cookie.Value == "" {
			if utils.IfArrayContains(cookieExemptPaths, path) {
				return next(c)
			}
			return c.Redirect(http.StatusSeeOther, "/app/form/api-key/")
		} else {
			keyValid := storage.IsKeyValid(cookie.Value)
			if !keyValid {
				ResetCookie(c)
				return c.Redirect(http.StatusSeeOther, "/app/form/api-key/")
			}
		}
		return next(c)
	}
}
