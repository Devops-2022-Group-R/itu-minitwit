package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlUsername := c.Param("username")
	author, err := userRepository.GetByUsername(urlUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if author == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "the username provided in the URL does not exist"})
		return
	}

	authUser := c.MustGet("user").(*models.User)
	if authUser.Username != urlUsername {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "the URL username did not match the Authorization header username"})
		return
	}

	var followTargetUserId int64
	if len(body.Follow) > 0 {
		followTarget, err := userRepository.GetByUsername(body.Follow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if followTarget == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "the username to follow does not exist"})
			return
		}
		followTargetUserId = followTarget.UserId
	}

	var unfollowTargetUserId int64
	if len(body.Unfollow) > 0 {
		unfollowTarget, err := userRepository.GetByUsername(body.Unfollow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if unfollowTarget == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "the username to unfollow does not exist"})
			return
		}
		unfollowTargetUserId = unfollowTarget.UserId
	}

	if len(body.Follow) > 0 {
		if err = userRepository.Follow(author.UserId, followTargetUserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(body.Unfollow) > 0 {
		if err = userRepository.Unfollow(author.UserId, unfollowTargetUserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if author == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "the username provided in the URL does not exist"})
		return
	}

	allFollowed, err := userRepository.AllFollowed(author.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
