package controllers

import (
	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"github.com/gin-gonic/gin"
)

func GetUserOrAdmin(c *gin.Context, userRepository database.IUserRepository) (*models.User, error) {
	urlUsername := c.Param("username")

	user := c.MustGet(UserKey).(*models.User)
	if c.MustGet(IsAdminKey).(bool) {
		var err error
		user, err = userRepository.GetByUsername(urlUsername)

		if err != nil {
			return nil, custom.NewInternalServerError(err)
		}

		if user == nil {
			return nil, custom.ErrUserNotFound
		}
	} else if user.Username != urlUsername {
		return nil, custom.ErrUrlUsernameNotMatchHeader
	}

	return user, nil
}
