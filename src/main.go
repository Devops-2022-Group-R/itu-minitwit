package main

import (
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

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
			initDb()
			return
		}
	}

	controllers.SetupRouter().Run()
}

// Creates the database tables.
func initDb() {
	db, err := database.ConnectDatabase(database.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	database.NewGormUserRepository(db).Migrate()
	database.NewGormMessageRepository(db).Migrate()
}
