package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
)

type FollowRequestBody struct {
	Follow string `form:"follow" json:"follow" binding:"required"`
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
		c.JSON(404, gin.H{"error": "the username provided in the URL does not exist"})
		return
	}

	// TODO: check if requestee = username in url. Or ignore requestee and use value from auth
	// TODO: forbidden if not logged in

	target, err := userRepository.GetByUsername(body.Follow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if target == nil {
		c.JSON(404, gin.H{"error": "the username to follow does not exist"})
		return
	}

	userRepository.Follow(author.UserId, target.UserId)

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
		c.JSON(404, gin.H{"error": "the username provided in the URL does not exist"})
		return
	}

	allFollowed, err := userRepository.AllFollowed(author.UserId)

	c.JSON(http.StatusOK, allFollowed)
}

// // Removes the current user as follower of the given user.
// func UnfollowController(c *gin.Context) {
// 	userRepository := c.MustGet("userRepository").(database.IUserRepository)

// 	username := c.Param("username")
// 	whom, err := userRepository.GetByUsername(username)
// 	if err != nil {
// 		c.JSON(404, nil)
// 		return
// 	}

// 	who, isLoggedIn := c.Get("user")

// 	if !isLoggedIn {
// 		c.JSON(401, nil)
// 		return
// 	}

// 	userRepository.Unfollow(who.(models.User).UserId, whom.UserId)

// 	flash(c, fmt.Sprintf("You are no longer following %s", username))

// 	c.Redirect(302, controllers.TimeLineUrl)
// }
