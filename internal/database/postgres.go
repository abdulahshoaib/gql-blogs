package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	dsn := "host=localhost user=trackme password=trackme123 dbname=trackme port=5432 sslmode=disable"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
