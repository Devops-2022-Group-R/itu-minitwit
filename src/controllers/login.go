package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

// Logs the user in.
func LoginPost(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Couldn't authenticate"})
		return
	}
	users := database.QueryDb(db, "select * from user where username = ?", username)

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Invalid username"})
		return
	} else if !pwdHash.CheckPasswordHash(password, users[0]["pw_hash"].(string)) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "password is incorrect"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func IsAuthenticated(c *gin.Context) bool {
	db := c.MustGet("db").(*sql.DB)
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		return false
	}
	users := database.QueryDb(db, "select * from user where username = ?", username)

	if len(users) == 0 || !pwdHash.CheckPasswordHash(password, users[0]["pw_hash"].(string)) {
		return false
	}
	return true
}
