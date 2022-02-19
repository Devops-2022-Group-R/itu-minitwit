package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase(databasePath string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return database, nil
}
