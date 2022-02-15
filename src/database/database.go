package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("./minitwit.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return database, nil
}
