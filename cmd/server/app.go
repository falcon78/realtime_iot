package main

import "gorm.io/gorm"

type app struct {
	db *gorm.DB
}

func newApp(db *gorm.DB) *app {
	return &app{
		db: db,
	}
}
