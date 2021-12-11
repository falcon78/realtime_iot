package main

import (
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"github.com/falcon78/realtime_iot/pkg/repository"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *app) handleWebsocket(c echo.Context) error {
	accessKey := c.Param("accessKey")

	channelRepo := repository.NewRepository(a.db)
	if _, err := channelRepo.GetChannelByAccessKey(accessKey); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Specified channel not found")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &realtime_update.Client{Conn: conn, Send: make(chan *realtime_update.Payload)}

	registerPayload := realtime_update.SubscriptionType{
		ChannelName: accessKey,
		Client:      client,
	}
	a.hub.Register <- &registerPayload

	go client.WritePump(a.hub, accessKey)
	return nil
}
