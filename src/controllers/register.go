package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

type RegisterRequestBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func RegisterController(c *gin.Context) {
	var body RegisterRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var errMsg string
	if body.Username == "" {
		errMsg = "username must not be empty"
	} else if body.Email == "" || !strings.Contains(body.Email, "@") {
		errMsg = "email address was not valid"
	} else if body.Password == "" {
		errMsg = "password must not be empty"
	}
	if errMsg != "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errMsg})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	if database.GetUserId(body.Username, db) != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "the username is already taken"})
		return
	}

	_, err := db.Exec(
		"insert into user (username, email, pw_hash) values (?, ?, ?)",
		body.Username,
		body.Email,
		pwdHash.GeneratePasswordHash(body.Password),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
