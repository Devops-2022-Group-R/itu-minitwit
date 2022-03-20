package controllers

import (
	"net/http"
	"strconv"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	"github.com/gin-gonic/gin"
)

func LatestController(c *gin.Context) {
	latestRepository := c.MustGet(LatestRepositoryKey).(database.ILatestRepository)

	latest, err := latestRepository.GetCurrent()
	if err != nil {
		custom.AbortWithError(c, custom.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"latest": latest})
}

func UpdateLatestMiddleware(c *gin.Context) {
	if values, ok := c.Request.URL.Query()["latest"]; ok {
		latestRepository := c.MustGet(LatestRepositoryKey).(database.ILatestRepository)
		newLatest, err := strconv.Atoi(values[0])
		if err == nil {
			latestRepository.Set(newLatest)
		}
	}

	c.Next()
}
