package main

import (
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"os"
	"time"

	"github.com/falcon78/realtime_iot/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	// don't use basic auth for this route
	e.POST("/api/record", app.postRecord, bodyDump())

	// websocket
	e.GET("/ws/:accessKey", app.handleWebsocket)

	// Serve static assets for frontend
	e.Static("/assets", "../../static/assets")
	e.File("/*", "../../static/index.html")

	timer := time.NewTicker(time.Second)

	dummyBroadcast := func() {
		for {
			select {
			case <-timer.C:
				message := realtime_update.Message{
					AccessKey: "channel",
					Payload: &realtime_update.Payload{
						ChannelOne:   0,
						ChannelTwo:   0,
						ChannelThree: 0,
						ChannelFour:  0,
					},
				}
				app.hub.Broadcast <- &message
			}
		}
	}

	go dummyBroadcast()
	go app.hub.Listen()

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
