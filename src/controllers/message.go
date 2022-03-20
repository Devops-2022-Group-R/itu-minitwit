package controllers

import (
	"net/http"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
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
		internal.AbortWithError(c, internal.NewInternalServerError(err))
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
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	if user == nil {
		internal.AbortWithError(c, internal.ErrUserNotFound)
		return
	}

	messages, err := messageRepository.GetByUserId(user.UserId, perPage)
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetFeedMessages(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)

	user := c.MustGet(UserKey).(*models.User)

	messages, err := messageRepository.GetByUserAndItsFollowers(user.UserId, perPage)
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, messages)
}

// PostUserMessage posts a non-empty message, with the current UTC time.
func PostUserMessage(c *gin.Context) {
	var body PostUserMessageRequestBody

	if err := c.BindJSON(&body); err != nil {
		internal.AbortWithError(c, internal.NewBadRequestErrorFromError(err))
		return
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	user, err := GetUserOrAdmin(c, userRepository)
	if err != (internal.HttpError{}) {
		internal.AbortWithError(c, err)
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
