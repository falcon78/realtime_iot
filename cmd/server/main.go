package main

import (
	"crypto/subtle"
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"os"
	"strings"
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
	//e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	//	fmt.Println(string(reqBody))
	//	fmt.Println(string(resBody))
	//}))
	protected := e.Group("")
	protected.Use(middleware.Gzip())

	user := strings.TrimSpace(os.Getenv("BASIC_AUTH_USER"))
	if user == "" || len(user) < 6 {
		panic("basic auth user name must be longer than 6 characters")
	}
	pass := strings.TrimSpace(os.Getenv("BASIC_AUTH_PASS"))
	if pass == "" || len(pass) < 6 {
		panic("basic auth password must be longer than 6 characters")
	}
	protected.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Routes
	protected.GET("/api/channels", app.getChannels)
	protected.POST("/api/channel/create/:channelName", app.createChannel)
	protected.DELETE("/api/channel/delete/:channelId", app.deleteChannel)
	protected.GET("/api/records/:channelId", app.getLatestRecords)
	e.POST("/api/record", app.postRecord)

	// Websocket
	protected.GET("/ws/:accessKey", app.handleWebsocket)

	// Serve static assets for frontend
	protected.Static("/assets", "../../static/assets")
	e.File("*", "../../static/index.html")

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
