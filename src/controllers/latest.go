package controllers

import (
	"net/http"
	"strconv"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
)

func LatestController(c *gin.Context) {
	latestRepository := c.MustGet(LatestRepositoryKey).(database.ILatestRepository)

	latest, err := latestRepository.GetCurrent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
