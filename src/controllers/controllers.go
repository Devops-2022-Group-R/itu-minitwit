package controllers

import (
	"log"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupRouter(openDatabase database.OpenDatabaseFunc) *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(beforeRequest(openDatabase))
	r.Use(UpdateLatestMiddleware)

	r.GET("/fllws/:username", FollowGetController)
	r.GET("/msgs", GetMessages)
	r.GET("/msgs/:username", GetUserMessages)
	r.POST("/register", RegisterController)
	r.GET("/latest", LatestController)

	authed := r.Group("/")
	authed.Use(AuthRequired())
	authed.GET("/login", LoginGet)
	authed.GET("/feed", GetFeedMessages)
	authed.POST("/fllws/:username", FollowPostController)
	authed.POST("/msgs/:username", PostUserMessage)

	return r
}

// Make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(openDatabase database.OpenDatabaseFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		gormDb, err := database.ConnectDatabase(openDatabase)
		if err != nil {
			log.Fatal(err)
		}

		c.Set(UserRepositoryKey, database.NewGormUserRepository(gormDb))
		c.Set(MessageRepositoryKey, database.NewGormMessageRepository(gormDb))
		c.Set(LatestRepositoryKey, database.NewGormLatestRepository(gormDb))

		c.Next()
	}
}
