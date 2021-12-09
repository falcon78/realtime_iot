package main

import (
	"fmt"
	"net/http"

	"github.com/falcon78/realtime_iot/pkg/repository"
	"github.com/labstack/echo/v4"
)

func (a *app) getRecords(c echo.Context) error {
	// user := new(models.UserParams)
	// if err := c.Bind(user); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Bad user paramaters")
	// }
	// user.UserName = strings.TrimSpace(user.UserName)
	// user.Password = strings.TrimSpace(user.Password)

	// if len(user.UserName) < 6 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Username should 6 letters or longer")
	// }
	// if len(user.UserName) < 8 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Password should 8 letters or longer")
	// }

	// add to db
	repo := repository.NewRepository(a.db)
	channels, err := repo.GetChannels()
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while fetching user",
		)
	}
	records, err := repo.GetAllRecordsByChannelId((*channels)[0].Id)

	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error occurred while fetching user",
		)
	}

	return c.JSON(http.StatusOK, records)
}
