package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"github.com/falcon78/realtime_iot/pkg/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func (a *app) getLatestRecords(c echo.Context) error {
	channelId, _ := strconv.Atoi(c.Param("channelId"))
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	repo := repository.NewRepository(a.db)
	records, err := repo.GetAllRecordsByChannelId(channelId, limit)
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

	current := time.Now()

	repo := repository.NewRepository(a.db)
	err := repo.AddRecord(r.AccessKey, r.ChannelOne, r.ChannelTwo, r.ChannelThree, r.ChannelFour, current)
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
			Timestamp:    current,
		},
	}

	return c.JSON(http.StatusOK, "Added record successfully")
}

func (a *app) downloadRecordCsv(c echo.Context) error {
	repo := repository.NewRepository(a.db)
	records, err := repo.GetAllRecordsByChannelKey(c.Param("channelKey"))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Channel with specified channel key not found",
		)
	} else if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while fetching records from db",
		)
	}

	csvContents := [][]string{
		{"timestamp", "channel1", "channel2", "channel3", "channel4"},
	}

	for _, r := range records {
		csvContents = append(csvContents, []string{
			r.Timestamp.String(),
			strconv.FormatFloat(r.ChannelOne, 'f', -1, 64),
			strconv.FormatFloat(r.ChannelTwo, 'f', -1, 64),
			strconv.FormatFloat(r.ChannelThree, 'f', -1, 64),
			strconv.FormatFloat(r.ChannelFour, 'f', -1, 64),
		})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/csv")

	writer := csv.NewWriter(c.Response().Writer)
	if err := writer.WriteAll(csvContents); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while writing data to csv",
		)
	}

	c.Response().Flush()
	return writer.Error()
}
