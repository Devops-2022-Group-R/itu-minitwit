package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
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

	if !strings.Contains(body.Email, "@") {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "email address was not valid"})
		return
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)

	if user, err := userRepository.GetByUsername(body.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "the username is already taken"})
		return
	}

	err := userRepository.Create(models.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: pwdHash.GeneratePasswordHash(body.Password),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
