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

	r.Use(ErrorHandleMiddleware())
	r.Use(CORSMiddleware())
	r.Use(beforeRequest(openDatabase))
	r.Use(UpdateLatestMiddleware)

	r.GET("/fllws/:username", FollowGetController)
	r.GET("/msgs", GetMessages)
	r.GET("/msgs/:username", GetUserMessages)
	r.POST("/register", RegisterController)
	r.GET("/latest", LatestController)

	r.PUT("/flag_tool/:msgid", FlagMessageById)
	r.GET("/flag_tool/msgs", GetAllMessages)

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

type returnedErr struct {
	Err        string `json:"error"`
	RelatedErr error  `json:"related_error"`
	Code       int    `json:"code"`
}

func ErrorHandleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		responseCode := 0

		if len(c.Errors) >= 1 {
			errors := make([]returnedErr, 0)

			for _, err := range c.Errors {
				switch err.Err.(type) {
				case HttpError:
					httpErr := err.Err.(HttpError)
					log.Printf("Http error (%d): %s, %s\n", httpErr.StatusCode, httpErr.Message, httpErr.RelatedErr)

					if !httpErr.Hidden {
						if httpErr.StatusCode > responseCode {
							responseCode = httpErr.StatusCode
						}

						errors = append(errors, returnedErr{httpErr.Message, httpErr.RelatedErr, httpErr.StatusCode})
					}
				default:
					log.Println("Internal server error: ", err)
					responseCode = 500
					errors = append(errors, returnedErr{"Internal server error", nil, 500})
				}

			}

			c.JSON(responseCode, gin.H{
				"errors": errors,
			})
		}
	}
}
