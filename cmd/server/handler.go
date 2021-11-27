package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app) home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Snippetbox")
}
