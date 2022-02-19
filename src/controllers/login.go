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
			c.JSON(http.StatusNotFound, gin.H{"status": "Invalid username"})
		} else if !pwdHash.CheckPasswordHash(body.Password, users[0]["pw_hash"].(string)) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "password is incorrect"})
		} else {
			//TODO: add login logic instead of session - gin authentication?
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}
}

// CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
// ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
// HEADERS = {'Connection': 'close',
//            'Content-Type': 'application/json',
//            f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()

	if hasAuth && user == "testuser" && password == "testpassword" {
		log.WithFields(log.Fields{
			"user": user,
		}).Info("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		return
	}
}
