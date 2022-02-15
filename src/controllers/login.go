package controllers

import (
	"database/sql"
	_ "log"
	"net/http"

	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	. "github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `form:"username" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Logs the user in.
func LoginPost(c *gin.Context) {
	var body Login
	if c.BindJSON(&body) == nil {
		db := c.MustGet("db").(*sql.DB)
		users := QueryDb(db, "select * from user where username = ?", body.Username)

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid username"})
		} else if !pwdHash.CheckPasswordHash(body.Password, users[0]["pw_hash"].(string)) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "password is incorrect"})
		} else {
			//TODO: add login logic instead of session - gin authentication?
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}
}
