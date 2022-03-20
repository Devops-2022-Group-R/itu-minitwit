package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/monitoring"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

type RegisterRequestBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"pwd" json:"pwd" binding:"required"`
}

var (
	ErrInvalidEmail = custom.NewHttpError(http.StatusUnprocessableEntity, "email address is not valid")
)

func RegisterController(c *gin.Context) {
	var body RegisterRequestBody

	if err := c.BindJSON(&body); err != nil {
		custom.AbortWithError(c, custom.NewBadRequestErrorFromError(err))
		return
	}

	if !strings.Contains(body.Email, "@") {
		custom.AbortWithError(c, ErrInvalidEmail)
		return
	}

	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)

	if user, err := userRepository.GetByUsername(body.Username); err != nil {
		custom.AbortWithError(c, custom.NewInternalServerError(err))
		return
	} else if user != nil {
		// We changed this from Conflict to Bad Request because the simulator
		// expects error code 400
		custom.AbortWithError(c, custom.NewBadRequestError("the username is already taken"))
		return
	}

	err := userRepository.Create(models.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: pwdHash.GeneratePasswordHash(body.Password),
	})

	if err != nil {
		custom.AbortWithError(c, custom.NewInternalServerError(err))
		return
	}

	monitoring.UserCount.Inc()

	c.JSON(http.StatusNoContent, nil)
}
