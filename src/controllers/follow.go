package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
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
		internal.AbortWithError(c, internal.NewBadRequestErrorFromError(err))
		return
	}

	user, err := GetUserOrAdmin(c, userRepository)
	if err != (internal.HttpError{}) {
		internal.AbortWithError(c, err)
		return
	}

	var followTargetUserId int64
	if len(body.Follow) > 0 {
		followTarget, err := userRepository.GetByUsername(body.Follow)
		if err != nil {
			internal.AbortWithError(c, internal.NewInternalServerError(err))
			return
		} else if followTarget == nil {
			internal.AbortWithError(c, internal.ErrUserNotFound)
			return
		}
		followTargetUserId = followTarget.UserId
	}

	var unfollowTargetUserId int64
	if len(body.Unfollow) > 0 {
		unfollowTarget, err := userRepository.GetByUsername(body.Unfollow)
		if err != nil {
			internal.AbortWithError(c, internal.NewInternalServerError(err))
			return
		} else if unfollowTarget == nil {
			internal.AbortWithError(c, internal.ErrUserNotFound)
			return
		}
		unfollowTargetUserId = unfollowTarget.UserId
	}

	if len(body.Follow) > 0 {
		if err := userRepository.Follow(user.UserId, followTargetUserId); err != nil {
			internal.AbortWithError(c, internal.NewInternalServerError(err))
			return
		}
	}

	if len(body.Unfollow) > 0 {
		if err := userRepository.Unfollow(user.UserId, unfollowTargetUserId); err != nil {
			internal.AbortWithError(c, internal.NewInternalServerError(err))
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
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	} else if author == nil {
		internal.AbortWithError(c, internal.ErrUserNotFound)
		return
	}

	allFollowed, err := userRepository.AllFollowed(author.UserId)
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
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
