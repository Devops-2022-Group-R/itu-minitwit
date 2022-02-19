package controllers

import (
	"net/http"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"github.com/gin-gonic/gin"
)

type MessageRequestBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Content  string `form:"content" json:"content" binding:"required"`
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
	var body MessageRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Content == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Message may not be empty"})
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

/*
// Registers a new message for the user.
func addMessage(c *gin.Context) {
	user, userLoggedIn := c.Get("user")

	if !userLoggedIn {
		c.JSON(401, nil)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	form := c.Request.Form
	text := form.Get("text")

	if text != "" {
		messageRepository := c.MustGet(messageRepositoryKey).(database.IMessageRepository)
		messageRepository.Create(models.Message{
			Author:  user.(models.User),
			Text:    text,
			PubDate: time.Now().Unix(),
		})

		flash(c, "Your message was recorded")
	}

	c.Redirect(302, timeLineUrl)
}
*/
