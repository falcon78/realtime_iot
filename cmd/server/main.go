package main

import (
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"github.com/falcon78/realtime_iot/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
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

	app := newApp(db, realtime_update.New())

	e := echo.New()
	e.Use()
	e.Use(middleware.Gzip())

	// Api Routes
	e.GET("/api/channels", app.getChannels, basicAuth())
	e.POST("/api/channel/create/:channelName", app.createChannel, basicAuth())
	e.DELETE("/api/channel/delete/:channelId", app.deleteChannel, basicAuth())
	e.GET("/api/records/:channelId", app.getLatestRecords, basicAuth())
	e.GET("/api/records/csv/:channelKey", app.downloadRecordCsv, basicAuth())
	// don't use basic auth for this route because channel
	// access key already acts like an authentication token
	e.POST("/api/record", app.postRecord)

	// websocket
	e.GET("/ws/:accessKey", app.handleWebsocket)

	// Serve static assets for frontend
	e.Static("/assets", "../../static/assets")
	e.File("/*", "../../static/index.html", basicAuth())

	// Start Websocket data broadcasting
	go app.hub.Listen()

	// Start http server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
