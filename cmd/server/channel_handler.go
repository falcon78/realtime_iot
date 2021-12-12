package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/falcon78/realtime_iot/pkg/repository"
	"github.com/labstack/echo/v4"
)

func (a *app) getChannels(c echo.Context) error {
	repo := repository.NewRepository(a.db)
	channels, err := repo.GetChannels()
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while fetching user",
		)
	}

	return c.JSON(http.StatusOK, channels)
}

func (a *app) createChannel(c echo.Context) error {
	channelName := strings.TrimSpace(c.Param("channelName"))
	if len(channelName) == 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"channel name must be longer than 1 character",
		)
	}

	repo := repository.NewRepository(a.db)
	if err := repo.CreateChannel(channelName); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while creating channel",
		)
	}

	return c.String(http.StatusOK, "Channel successfully created")
}

func (a *app) deleteChannel(c echo.Context) error {
	channelId, _ := strconv.Atoi(c.Param("channelId"))

	repo := repository.NewRepository(a.db)
	if err := repo.DeleteChannel(channelId); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while deleting channel",
		)
	}

	return c.String(http.StatusOK, "Channel successfully deleted")
}
