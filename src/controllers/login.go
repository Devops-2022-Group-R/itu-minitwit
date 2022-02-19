package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

// Logs the user in.
func LoginGet(c *gin.Context) {

	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Couldn't authenticate"})
		return
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	user, err := userRepository.GetByUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Invalid username"})
		return
	} else if !pwdHash.CheckPasswordHash(password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "password is incorrect"})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

func IsAuthenticated(c *gin.Context) bool {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		return false
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	user, err := userRepository.GetByUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	if user == nil || !pwdHash.CheckPasswordHash(password, user.PasswordHash) {
		return false
	}
	return true
}
