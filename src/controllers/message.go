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

// GetMessages returns the latest messages, limited by the number of messages per page.
func GetMessages(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)
	messages, err := messageRepository.GetWithLimit(perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetUserMessages returns the latest messages by the user, limited by the number of messages per page.
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

func GetFeedMessages(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)

	user := c.MustGet(UserKey).(*models.User)

	messages, err := messageRepository.GetByUserAndItsFollowers(user.UserId, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// PostUserMessage posts a non-empty message, with the current UTC time.
func PostUserMessage(c *gin.Context) {
	var body PostUserMessageRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlUsername := c.Param("username")

	user := c.MustGet(UserKey).(*models.User)
	if user.Username == "simulator" {
		userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
		var err error
		user, err = userRepository.GetByUsername(urlUsername)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
	} else if user.Username != urlUsername {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "the URL username did not match the Authorization header username"})
		return
	}

	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)
	messageRepository.Create(models.Message{
		Author:  *user,
		Text:    body.Content,
		PubDate: time.Now().UTC().Unix(),
		Flagged: false,
	})

	c.JSON(http.StatusNoContent, nil)
}
