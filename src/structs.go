package main

import "github.com/gin-gonic/gin"

type User struct {
	UserId       int64
	Username     string
	Email        string
	PasswordHash string
}

type Message struct {
	Email    string // Email of the user submitting this message
	Username string // Username of the user submitting this message
	PubDate  int64  // The publish timestamp as UNIX
	Text     string // The message itself
}

type LayoutData struct {
	Context *gin.Context
	User    User // Me
}

type DataProvider interface {
	setContext(ctx *gin.Context)
	setUser(user User)
}

func (ld LayoutData) setContext(ctx *gin.Context) {
	ld.Context = ctx
}

func (ld LayoutData) setUser(user User) {
	ld.User = user
}

type TimelineData struct {
	LayoutData

	IsPublicTimeline bool
	IsMyTimeline     bool // Used if IsPublicTimeline is false
	IsFollowed       bool // Used if IsMyTimeline is false
	HasMessages      bool

	ProfileUser User // Used if IsMyTimeline is false, represents the user of the profile you visit

	Messages []Message
}

type LoginData struct {
	LayoutData

	Username string
	ErrorMsg string
}

type RegisterData struct {
	LayoutData

	Username string
	Email    string
	ErrorMsg string
}
