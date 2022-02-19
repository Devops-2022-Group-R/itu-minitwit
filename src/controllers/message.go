package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageRequestBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Message  string `form:"text" json:"text" binding:"required"`
}

func GetMessage(c *gin.Context) {
}

// Look into authorization
func PostMessage(c *gin.Context) {
	var body MessageRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Message == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Message may not be empty"})
		return
	}

	// db := c.MustGet("db").(*sql.DB)
	// database.QueryDb(db, "insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, 0)",
	// 	user.(User).UserId, body.Message, time.Now().UTC().Unix())

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
