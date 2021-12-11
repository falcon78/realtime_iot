package main

import (
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"strings"
)

func basicAuth() echo.MiddlewareFunc {
	user := strings.TrimSpace(os.Getenv("BASIC_AUTH_USER"))
	if user == "" || len(user) < 6 {
		panic("basic auth user name must be longer than 6 characters")
	}
	pass := strings.TrimSpace(os.Getenv("BASIC_AUTH_PASS"))
	if pass == "" || len(pass) < 6 {
		panic("basic auth password must be longer than 6 characters")
	}

	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
func bodyDump() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Println(string(reqBody))
		fmt.Println(string(resBody))
	})
}
