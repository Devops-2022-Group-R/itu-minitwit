package database

import (
	"log"

	"gorm.io/gorm"
)

const DatabasePath = "./minitwit.db"

type OpenDatabaseFunc = func() gorm.Dialector

func ConnectDatabase(openDatabase OpenDatabaseFunc) (*gorm.DB, error) {
	database, err := gorm.Open(openDatabase(), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return database, nil
}

func InitDatabase(openDatabase OpenDatabaseFunc) {
	db, err := ConnectDatabase(openDatabase)
	if err != nil {
		log.Fatal(err)
	}

	NewGormUserRepository(db).Migrate()
	NewGormMessageRepository(db).Migrate()
}
