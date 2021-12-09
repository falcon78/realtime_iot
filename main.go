package main

import (
	"crypto/subtle"
	"os"

	"github.com/falcon78/realtime_iot/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var MigrationDir = "file://migrations"

func main() {
	db, err := utils.GetDb()
	if err != nil {
		panic(err)
	}

	// Run database migrations
	if m, err := migrate.New(MigrationDir, utils.GetDatabaseStringForMigrate()); err != nil {
		panic(err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
	}

	app := newApp(db)

	e := echo.New()
	e.Use(middleware.Gzip())

	// TODO: import basic auth credentials from env
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("falcon")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Routes
	e.POST("/api/getRecords", app.getRecords)

	// Serve static assets for frontend
	e.Static("/assets", "static/assets")
	e.File("*", "static/index.html")

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
