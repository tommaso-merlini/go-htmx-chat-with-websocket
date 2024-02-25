package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

	"roomate/pkg/sb"
	"roomate/types"
)

func WithUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().URL.Path, "/public") {
			return next(c)
		}
		accessToken, err := getAccessToken(c)
		if err != nil {
			return next(c)
		}
		resp, err := sb.Client.Auth.User(c.Request().Context(), accessToken)
		if err != nil {
			return next(c)
		}
		user := types.AuthenticatedUser{
			AuthID:     resp.ID,
			Email:      resp.Email,
			IsLoggedIn: true,
		}
		c.Set(types.UserContextKey, user)
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
	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return "", errors.New("session not found")
	}
	return session.Values[sessionAccessTokenKey].(string), nil
}
