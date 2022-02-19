package controllers

import (
	"net/http"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"github.com/gin-gonic/gin"
)

type PostUserMessageRequestBody struct {
	Content string `form:"content" json:"content" binding:"required"`
}

func GetMessages(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)
	messages, err := messageRepository.GetWithLimit(perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetUserMessages(c *gin.Context) {
	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)

	user, err := userRepository.GetByUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	messages, err := messageRepository.GetByUserId(user.UserId, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func PostUserMessage(c *gin.Context) {
	var body PostUserMessageRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)

	user, err := userRepository.GetByUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	messageRepository.Create(models.Message{
		Author:  *user,
		Text:    body.Content,
		PubDate: time.Now().UTC().Unix(),
		Flagged: false,
	})

	c.JSON(http.StatusNoContent, nil)
}
