package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DatabasePath     = "./minitwit.db"
	TestDatabasePath = "./minitwit-test.db"
)

func ConnectDatabase(databasePath string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return database, nil
}

func InitDatabase(databasePath string) {
	db, err := ConnectDatabase(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	NewGormUserRepository(db).Migrate()
	NewGormMessageRepository(db).Migrate()
}

func NukeDatabase(databasePath string) {
	if err := os.Remove(databasePath); err != nil {
		log.Fatal(err)
	}
}
