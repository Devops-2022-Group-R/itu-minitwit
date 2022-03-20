package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/custom"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
)

type FollowRequestBody struct {
	Follow   string `form:"follow" json:"follow"`
	Unfollow string `form:"unfollow" json:"unfollow"`
}

// Adds the current user as follower of the given user.
func FollowPostController(c *gin.Context) {
	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)

	var body FollowRequestBody
	if err := c.BindJSON(&body); err != nil {
		custom.AbortWithError(c, custom.NewBadRequestErrorFromError(err))
		return
	}

	user, err := GetUserOrAdmin(c, userRepository)
	if err != nil {
		custom.AbortWithError(c, err)
		return
	}

	var followTargetUserId int64
	if len(body.Follow) > 0 {
		followTarget, err := userRepository.GetByUsername(body.Follow)
		if err != nil {
			custom.AbortWithError(c, custom.NewInternalServerError(err))
			return
		} else if followTarget == nil {
			custom.AbortWithError(c, custom.ErrUserNotFound)
			return
		}
		followTargetUserId = followTarget.UserId
	}

	var unfollowTargetUserId int64
	if len(body.Unfollow) > 0 {
		unfollowTarget, err := userRepository.GetByUsername(body.Unfollow)
		if err != nil {
			custom.AbortWithError(c, custom.NewInternalServerError(err))
			return
		} else if unfollowTarget == nil {
			custom.AbortWithError(c, custom.ErrUserNotFound)
			return
		}
		unfollowTargetUserId = unfollowTarget.UserId
	}

	if len(body.Follow) > 0 {
		if err := userRepository.Follow(user.UserId, followTargetUserId); err != nil {
			custom.AbortWithError(c, custom.NewInternalServerError(err))
			return
		}
	}

	if len(body.Unfollow) > 0 {
		if err := userRepository.Unfollow(user.UserId, unfollowTargetUserId); err != nil {
			custom.AbortWithError(c, custom.NewInternalServerError(err))
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

func FollowGetController(c *gin.Context) {
	userRepository := c.MustGet(UserRepositoryKey).(database.IUserRepository)

	urlUsername := c.Param("username")
	author, err := userRepository.GetByUsername(urlUsername)
	if err != nil {
		custom.AbortWithError(c, custom.NewInternalServerError(err))
		return
	} else if author == nil {
		custom.AbortWithError(c, custom.ErrUserNotFound)
		return
	}

	allFollowed, err := userRepository.AllFollowed(author.UserId)
	if err != nil {
		custom.AbortWithError(c, custom.NewInternalServerError(err))
		return
	}

	usernames := make([]string, len(allFollowed))
	for i, user := range allFollowed {
		usernames[i] = user.Username
	}

	c.JSON(http.StatusOK, gin.H{
		"follows": usernames,
	})
}
