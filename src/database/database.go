package database

import (
	"log"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
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

	userRepository := NewGormUserRepository(db)
	userRepository.Migrate()
	NewGormMessageRepository(db).Migrate()
	NewGormLatestRepository(db).Migrate()

	// The simulator needs to be a default user
	if err = userRepository.Create(models.User{
		Username:     "simulator",
		Email:        "unused@email.rip",
		PasswordHash: "pbkdf2:sha256:260000$MYtvPfLTCXY74kXG$44274a66ed9a6d08471cfc77337ae10497ba3c4a4cb4adec02227305ce663378",
	}); err != nil {
		log.Fatal(err)
	}
}
