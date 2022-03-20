package database

import (
	"log"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	"gorm.io/gorm"
)

type OpenDatabaseFunc = func() gorm.Dialector

var Db *gorm.DB

func ConnectDatabase(openDatabase OpenDatabaseFunc) (*gorm.DB, error) {
	if Db != nil {
		return Db, nil
	}

	database, err := gorm.Open(openDatabase(), &gorm.Config{
		Logger: internal.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDb, err := database.DB()
	if err != nil {
		return nil, err
	}

	// Arbitrarily picked numbers
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(50)

	// Microsoft sql server timeouts after 30 seconds
	sqlDb.SetConnMaxIdleTime(time.Second * 29)

	Db = database

	InitDatabase(openDatabase)

	return Db, nil
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

	user, err := userRepository.GetByUsername("simulator")
	if err != nil {
		log.Fatal(err)
	}

	if user == nil {
		// The simulator needs to be a default user
		if err = userRepository.Create(models.User{
			Username:     "simulator",
			Email:        "unused@email.rip",
			PasswordHash: pwdHash.GeneratePasswordHash("super_safe!"),
		}); err != nil {
			log.Fatal(err)
		}
	}
}
