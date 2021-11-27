package main

import (
	"os"

	"github.com/falcon78/realtime_iot/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
)

var MigrationDir = "file://../../migrations"

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

	// Routes
	e.GET("/", app.home)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
