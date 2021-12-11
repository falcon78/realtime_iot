package main

import (
	"github.com/falcon78/realtime_iot/pkg/realtime_update"
	"gorm.io/gorm"
)

type app struct {
	db  *gorm.DB
	hub *realtime_update.Hub
}

func newApp(db *gorm.DB, hub *realtime_update.Hub) *app {
	return &app{
		db:  db,
		hub: hub,
	}
}
