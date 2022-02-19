package controllers

import (
	"log"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
)

func SetupRouter(openDatabase database.OpenDatabaseFunc) *gin.Engine {
	r := gin.Default()

	r.Use(beforeRequest(openDatabase))

	r.POST("/fllws/:username", FollowPostController)
	r.GET("/fllws/:username", FollowGetController)
	r.GET("/msgs", GetMessages)
	r.GET("/msgs/:username", GetUserMessages)
	r.POST("/msgs/:username", PostUserMessage)
	r.POST(LoginUrl, LoginPost)
	r.POST("/register", RegisterController)

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

		c.Next()
	}
}
