package repository

import (
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repo {
	return &Repo{
		DB: db,
	}
}
