package main

import (
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/monitoring"
	_ "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	"github.com/denisenkom/go-mssqldb/azuread"
)

func init() {
	monitoring.Initialise(openDatabase)
}

func main() {
	if len(os.Args) > 1 {
		input := os.Args[1]
		if strings.EqualFold("initDb", input) {
			database.InitDatabase(openDatabase)
			return
		}
	}

	err := controllers.SetupRouter(openDatabase).Run()

	if err != nil {
		custom.Logger.Fatalf("failed to start: %v", err)
	}
}

func openDatabase() gorm.Dialector {
	env := os.Getenv("ENVIRONMENT")
	if env == "PRODUCTION" {
		connString, exists := os.LookupEnv("SQLCONNSTR_CONNECTION_STRING")
		if !exists {
			custom.Logger.Fatalln("SQLCONNSTR_CONNECTION_STRING environment variable not set")
		}

		return sqlserver.New(sqlserver.Config{
			DSN:        connString,
			DriverName: azuread.DriverName,
		})
	} else {
		return sqlite.Open("minitwit.db")
	}
}
