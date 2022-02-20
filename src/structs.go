package main

import "github.com/Devops-2022-Group-R/itu-minitwit/src/models"

type LayoutData struct {
	Flashes []string
	User    models.User // Me
}

type DataProvider interface {
	initLayoutData()
	setFlashes([]string)
	setUser(user models.User)
}

func (ld *LayoutData) setFlashes(flashes []string) {
	ld.Flashes = flashes
}

func (ld *LayoutData) setUser(user models.User) {
	ld.User = user
}

type TimelineData struct {
	*LayoutData

	IsPublicTimeline bool
	IsMyTimeline     bool // Used if IsPublicTimeline is false
	IsFollowed       bool // Used if IsMyTimeline is false
	HasMessages      bool

	ProfileUser models.User // Used if IsMyTimeline is false, represents the user of the profile you visit

	Messages []models.Message
}

func (t *TimelineData) initLayoutData() {
	t.LayoutData = &LayoutData{}
}

type LoginData struct {
	*LayoutData

	Username string
	ErrorMsg string
}

func (t *LoginData) initLayoutData() {
	t.LayoutData = &LayoutData{}
}

type RegisterData struct {
	*LayoutData

	Username string
	Email    string
	ErrorMsg string
}

func (t *RegisterData) initLayoutData() {
	t.LayoutData = &LayoutData{}
}
