package main

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	_ "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

const (
	debug     = true
	secretKey = "development key"
)

type Row = map[string]interface{}

func main() {
	if debug {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	if len(os.Args) > 1 {
		input := os.Args[1]
		if strings.EqualFold("initDb", input) {
			database.InitDatabase(openDatabase)
			return
		}
	}

	controllers.SetupRouter(openDatabase).Run()
}

func openDatabase() gorm.Dialector {
	env := os.Getenv("ENVIRONMENT")
	if env == "PRODUCTION" {
		connString, exists := os.LookupEnv("CONNECTION_STRING")
		if !exists {
			log.Fatal("CONNECTION_STRING environment variable not set")
		}

		return postgres.Open(connString)
	} else {
		return sqlite.Open("minitwit.db")
	}
}
