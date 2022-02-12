package main

type User struct {
	UserId       int
	Username     string
	Email        string
	PasswordHash string
}

type Message struct {
	Email    string // Email of the user submitting this message
	Username string // Username of the user submitting this message
	PubDate  int    // FIXME: I think, depends on the database
	Text     string // The message itself
}

type LayoutData struct {
	User User // Me
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
	Username string
	ErrorMsg string
}

type RegisterData struct {
	Username string
	Email    string
	ErrorMsg string
}
