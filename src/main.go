package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", database.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
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

// Format a timestamp for display.
func formatDateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02 @ 15:04")
}

// Return the gravatar image for the given email address.
func gravatarUrl(email string, size int) string {
	email = strings.ToLower(strings.TrimSpace(email))

	hash := md5.New()
	hash.Write([]byte(email))

	hex := fmt.Sprintf("%x", hash.Sum(nil))

	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hex, size)
}
