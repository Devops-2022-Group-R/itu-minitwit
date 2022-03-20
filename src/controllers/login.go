package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

var (
	ErrInvalidUsername    = custom.NewHttpError(http.StatusNotFound, "invalid username")
	ErrIncorrectPassword  = custom.NewHttpError(http.StatusUnauthorized, "incorrect password")
	ErrMissingCredentials = custom.NewHttpError(http.StatusUnauthorized, "missing authentication credentials")
)

// Logs the user in.
// Essentially a test of the AuthRequired middleware.
func LoginGet(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
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

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUsername, err := GetAuthState(c)
		if authUsername == "" || err != nil {
			custom.AbortWithError(c, err)

			return
		}

		userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)
		user, err := userRepository.GetByUsername(authUsername)
		if err != nil {
			custom.AbortWithError(c, custom.NewInternalServerError(err))
			return
		}
		if user == nil {
			custom.AbortWithError(c, custom.ErrUserNotFound)
			return
		}

		if c.GetHeader("Authorization") == "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
			c.Set(IsAdminKey, true)
		} else {
			c.Set(IsAdminKey, false)
		}

		c.Set(UserKey, user)
		c.Next()
	}
}
