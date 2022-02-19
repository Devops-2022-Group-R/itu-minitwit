package controllers

/*
// Adds the current user as follower of the given user.
func followUser(c *gin.Context) {
	userRepository := c.MustGet("userRepository").(database.IUserRepository)

	username := c.Param("username")
	whom, err := userRepository.GetByUsername(username)
	if err != nil {
		c.JSON(404, nil)
		return
	}

	who, isLoggedIn := c.Get("user")

	if !isLoggedIn {
		c.JSON(401, nil)
		return
	}

	userRepository.Follow(who.(models.User).UserId, whom.UserId)

	flash(c, fmt.Sprintf("You are now following %s", username))

	c.Redirect(302, timeLineUrl)
}

// Removes the current user as follower of the given user.
func unfollowUser(c *gin.Context) {
	userRepository := c.MustGet("userRepository").(database.IUserRepository)

	username := c.Param("username")
	whom, err := userRepository.GetByUsername(username)
	if err != nil {
		c.JSON(404, nil)
		return
	}

	who, isLoggedIn := c.Get("user")

	if !isLoggedIn {
		c.JSON(401, nil)
		return
	}

	userRepository.Unfollow(who.(models.User).UserId, whom.UserId)

	flash(c, fmt.Sprintf("You are no longer following %s", username))

	c.Redirect(302, controllers.TimeLineUrl)
}
*/
