package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

var (
	ErrInvalidUsername    = errors.New("invalid username")
	ErrIncorrectPassword  = errors.New("incorrect password")
	ErrMissingCredentials = errors.New("missing authentication credentials")
)

// Logs the user in.
func LoginGet(c *gin.Context) {
	_, err := GetAuthState(c)

	switch err {
	case nil:
		c.JSON(http.StatusNoContent, nil)
	case ErrInvalidUsername:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case ErrIncorrectPassword, ErrMissingCredentials:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetAuthState(c *gin.Context) (string, error) {
	username, password, hasAuth := c.Request.BasicAuth()

	if !hasAuth {
		return username, ErrMissingCredentials
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
	user, err := userRepository.GetByUsername(username)
	if err != nil {
		return username, err
	}

	if user == nil {
		return username, ErrInvalidUsername
	}

	if !pwdHash.CheckPasswordHash(password, user.PasswordHash) {
		return username, ErrIncorrectPassword
	}

	return username, nil
}
