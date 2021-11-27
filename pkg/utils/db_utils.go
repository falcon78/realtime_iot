package utils

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseStringForMigrate() string {
	var sslMode string
	if os.Getenv("SSL_MODE") != "enable" {
		sslMode = "sslmode=disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		sslMode,
	)
}

func GetDatabaseString() string {
	var sslMode string
	if os.Getenv("SSL_MODE") != "enable" {
		sslMode = "sslmode=disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s %s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		sslMode,
	)
}

func GetDb() (db *gorm.DB, err error) {
	dsn := GetDatabaseString()
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}
