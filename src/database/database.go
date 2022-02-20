package database

import (
	"log"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	"gorm.io/gorm"
)

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
		PasswordHash: pwdHash.GeneratePasswordHash("super_safe!"),
	}); err != nil {
		log.Fatal(err)
	}
}
