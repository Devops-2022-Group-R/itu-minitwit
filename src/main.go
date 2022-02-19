package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	_ "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

var databasePath = "./minitwit.db"

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

	setupRouter().Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(secretKey))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(beforeRequest)
	r.Static("/static", "./src/static")

	//r.GET("/:username/follow", followUser)
	//r.GET("/:username/unfollow", unfollowUser)
	r.GET("/msgs", controllers.GetMessages)
	r.GET("/msgs/:username", controllers.GetUserMessages)
	r.POST("/msgs/:username", controllers.PostUserMessage)
	r.POST(controllers.LoginUrl, controllers.LoginPost)
	r.POST("/register", controllers.RegisterController)

	return r
}

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Creates the database tables.
func initDb() {
	db, err := database.ConnectDatabase(databasePath)
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

// Make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(c *gin.Context) {
	gormDb, err := database.ConnectDatabase(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	userRepository := database.NewGormUserRepository(gormDb)
	c.Set(controllers.UserRepositoryKey, database.NewGormUserRepository(gormDb))
	c.Set(controllers.MessageRepositoryKey, database.NewGormMessageRepository(gormDb))

	session := sessions.Default(c)
	if userId := session.Get("user_id"); userId != nil {
		user, err := userRepository.GetByID(userId.(int64))
		if err != nil {
			log.Fatal(err)
		}
		c.Set("user", user)
	}

	c.Next()
}
