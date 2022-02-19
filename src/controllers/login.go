package controllers

import (
	"database/sql"
	_ "log"
	"net/http"

	. "github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	"github.com/gin-gonic/gin"
)

// This is not needed when we use BasicAuth ?
type Login struct {
	Username string `form:"username" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
// ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
// HEADERS = {'Connection': 'close',
//            'Content-Type': 'application/json',
//            f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}


// Logs the user in.
func LoginPost(c *gin.Context) { // add bool return type?
	var body Login // remove this?
	if c.BindJSON(&body) == nil {
		db := c.MustGet("db").(*sql.DB)
		username, password, hasAuth := c.Request.BasicAuth()
		users := QueryDb(db, "select * from user where username = ?", username)

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "Invalid username"})
		} else if !pwdHash.CheckPasswordHash(password, users[0]["pw_hash"].(string)) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "password is incorrect"})
		} else if hasAuth {
			c.JSON(http.StatusNoContent, nil)
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Couldn't authenticate"}) // // c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.Abort()
			return
		}
	}
}