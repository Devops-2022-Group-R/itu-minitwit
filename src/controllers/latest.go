package controllers

import (
	"net/http"
	"strconv"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
	"github.com/gin-gonic/gin"
)

func LatestController(c *gin.Context) {
	latestRepository := c.MustGet(LatestRepositoryKey).(database.ILatestRepository)

	latest, err := latestRepository.GetCurrent()
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"latest": latest})
}

func UpdateLatestMiddleware(c *gin.Context) {
	if values, ok := c.Request.URL.Query()["latest"]; ok {
		latestRepository := c.MustGet(LatestRepositoryKey).(database.ILatestRepository)
		newLatest, err := strconv.Atoi(values[0])
		if err == nil {
			if err = latestRepository.Set(newLatest); err != nil {
				internal.Logger.Printf("ERROR - updating latest middleware failed: %v", err)
			}
		}
	}

	c.Next()
}
