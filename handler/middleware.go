package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"roomate/pkg/sb"
	"roomate/types"
)

func WithUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		print(0)
		if strings.Contains(c.Request().URL.Path, "/public") {
			return next(c)
		}

		print(1)
		accessToken, err := getAccessToken(c)
		if err != nil {
			return next(c)
		}

		print(2)
		resp, err := sb.Client.Auth.User(c.Request().Context(), accessToken)
		if err != nil {
			return next(c)
		}
		print(3)
		user := types.AuthenticatedUser{
			AuthID:     resp.ID,
			Email:      resp.Email,
			IsLoggedIn: true,
		}
		SetAuthenticatedUser(c, user)
		return next(c)
	}
}

func WithAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().URL.Path, "/public") {
			return next(c)
		}
		user := GetAuthenticatedUser(c)
		if !user.IsLoggedIn {
			path := c.Request().URL.Path
			http.Redirect(c.Response(), c.Request(), "/login?to="+path, http.StatusSeeOther)
			return nil
		}
		c.Set(types.UserContextKey, user)
		return next(c)
	}
}

func getAccessToken(c echo.Context) (string, error) {
	cookie, err := c.Cookie(sessionAccessTokenKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
