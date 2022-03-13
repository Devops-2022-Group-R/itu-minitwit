package controllers

import (
	"log"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log request
		internal.Logger.Printf("[GIN][%3d][%13v][IP: %15s][%-7s][%s] %s\n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}

func SetupRouter(openDatabase database.OpenDatabaseFunc) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(LoggingMiddleware())
	r.Use(CORSMiddleware())
	r.Use(beforeRequest(openDatabase))
	r.Use(UpdateLatestMiddleware)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
