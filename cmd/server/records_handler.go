package main

import (
	"fmt"
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"github.com/falcon78/realtime_iot/pkg/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (a *app) getLatestRecords(c echo.Context) error {
	channelId, _ := strconv.Atoi(c.Param("channelId"))

	repo := repository.NewRepository(a.db)
	records, err := repo.GetAllRecordsByChannelId(channelId)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error fetching records",
		)
	}

	return c.JSON(http.StatusOK, records)
}

type recordRequest struct {
	AccessKey    string  `json:"accessKey"`
	ChannelOne   float64 `json:"channelOne"`
	ChannelTwo   float64 `json:"channelTwo"`
	ChannelThree float64 `json:"channelThree"`
	ChannelFour  float64 `json:"channelFour"`
}

func (a *app) postRecord(c echo.Context) error {
	r := new(recordRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	repo := repository.NewRepository(a.db)
	err := repo.AddRecord(r.AccessKey, r.ChannelOne, r.ChannelTwo, r.ChannelThree, r.ChannelFour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	a.hub.Broadcast <- &realtime_update.Message{
		AccessKey: r.AccessKey,
		Payload: &realtime_update.Payload{
			ChannelOne:   r.ChannelOne,
			ChannelTwo:   r.ChannelTwo,
			ChannelThree: r.ChannelThree,
			ChannelFour:  r.ChannelFour,
		},
	}

	return c.JSON(http.StatusOK, "Added record successfully")
}
