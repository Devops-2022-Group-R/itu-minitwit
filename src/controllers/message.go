package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
)

const messagesPerPage = 30
const USERID = -42 // TODO: Figure out user ID

type MessageRequestBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Content  string `form:"content" json:"content" binding:"required"`
}

func GetMessages(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	query := `SELECT message.*, user.*
			  FROM message, user
			  WHERE message.flagged = 0
			      AND message.author_id = user.user_id
			  ORDER BY message.pub_date DESC
			  LIMIT ?`
	rows := database.QueryDb(db, query, messagesPerPage)

	println(rows) // To compile
}

func GetUserMessages(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	query := `SELECT message.*, user.*
			  FROM message, user 
			  WHERE message.flagged = 0
			  	  AND user.user_id = message.author_id
				  AND user.user_id = ?
			  ORDER BY message.pub_date DESC
			  LIMIT ?`
	rows := database.QueryDb(db, query, USERID, messagesPerPage)

	println(rows) // To compile
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

	query := "INSERT INTO message (author_id, text, pub_date, flagged) VALUES (?, ?, ?, 0)"
	db := c.MustGet("db").(*sql.DB)
	database.QueryDb(db, query, USERID, body.Content, time.Now().UTC().Unix())

	c.JSON(http.StatusNoContent, nil)
}
