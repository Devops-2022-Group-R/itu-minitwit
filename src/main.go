package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	_ "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	pwdHash "github.com/Devops-2022-Group-R/itu-minitwit/src/password"
)

const (
	perPage   = 30
	debug     = true
	secretKey = "development key"

	timeLineUrl       = "/"
	publicTimelineUrl = "/public"
	loginUrl          = "/login"

	flashesKey = "flashes"

	userRepositoryKey    = "userRepository"
	messageRepositoryKey = "messageRepository"
)

var databasePath = "./minitwit.db"

type Row = map[string]interface{}

func main() {
	if debug {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	if len(os.Args) > 1 {
		input := os.Args[1]
		if strings.EqualFold("initDb", input) {
			initDb()
			return
		}
	}

	setupRouter().Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(secretKey))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(beforeRequest)
	r.Static("/static", "./src/static")

	r.GET(timeLineUrl, timeline)
	r.GET(publicTimelineUrl, publicTimeline)
	r.GET("/:username", userTimeline)
	r.GET("/:username/follow", followUser)
	r.GET("/:username/unfollow", unfollowUser)
	r.GET("/msgs/:username", controllers.GetMessage)
	r.POST("/msgs/:username", controllers.PostMessage)
	r.GET(loginUrl, loginGet)
	r.POST(loginUrl, controllers.LoginPost)
	r.POST("/register", controllers.RegisterController)

	return r
}

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Creates the database tables.
func initDb() {
	db, err := database.ConnectDatabase(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	database.NewGormUserRepository(db).Migrate()
	database.NewGormMessageRepository(db).Migrate()
}

// Format a timestamp for display.
func formatDateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02 @ 15:04")
}

// Return the gravatar image for the given email address.
func gravatarUrl(email string, size int) string {
	email = strings.ToLower(strings.TrimSpace(email))

	hash := md5.New()
	hash.Write([]byte(email))

	hex := fmt.Sprintf("%x", hash.Sum(nil))

	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hex, size)
}

func flash(c *gin.Context, message string) {
	session := sessions.Default(c)

	var flashes []string
	if f := session.Get(flashesKey); f != nil {
		flashes = f.([]string)
	} else {
		flashes = make([]string, 1)
	}

	flashes = append(flashes, message)

	session.Set(flashesKey, flashes)
	session.Save()
}

func getFlashes(c *gin.Context) []string {
	session := sessions.Default(c)

	if f := session.Get(flashesKey); f != nil {
		session.Set(flashesKey, make([]string, 0))
		session.Save()

		return f.([]string)
	} else {
		return nil
	}
}

// Make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(c *gin.Context) {
	gormDb, err := database.ConnectDatabase(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	userRepository := database.NewGormUserRepository(gormDb)
	c.Set(userRepositoryKey, database.NewGormUserRepository(gormDb))
	c.Set(messageRepositoryKey, database.NewGormMessageRepository(gormDb))

	session := sessions.Default(c)
	if userId := session.Get("user_id"); userId != nil {
		user, err := userRepository.GetByID(userId.(int64))
		if err != nil {
			log.Fatal(err)
		}
		c.Set("user", user)
	}

	c.Next()
}

// Shows a users timeline or if no user is logged in it will redirect to the
// public timeline. This timeline shows the user's messages as well as all the
// messages of followed users.
func timeline(c *gin.Context) {
	messageRepository := c.MustGet(messageRepositoryKey).(database.IMessageRepository)

	var user models.User
	if u, isLoggedIn := c.Get("user"); isLoggedIn {
		user = u.(models.User)
	} else {
		c.Redirect(302, publicTimelineUrl)
		return
	}

	messages, err := messageRepository.GetByUserAndItsFollowers(user.UserId, perPage)
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(c, "timeline.html", &TimelineData{
		ProfileUser: user,

		IsPublicTimeline: false,
		IsMyTimeline:     true,
		IsFollowed:       false,
		HasMessages:      len(messages) > 0,
		Messages:         messages,
	})
}

// Displays the latest messages of all users.
func publicTimeline(c *gin.Context) {
	messageRepository := c.MustGet(messageRepositoryKey).(database.IMessageRepository)

	messages, err := messageRepository.GetWithLimit(perPage)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Could not fetch messages."})
		return
	}

	renderTemplate(c, "timeline.html", &TimelineData{
		IsPublicTimeline: true,
		IsMyTimeline:     false,
		IsFollowed:       true,
		HasMessages:      len(messages) > 0,
		Messages:         messages,
	})
}

// Display's a users tweets.
func userTimeline(c *gin.Context) {
	userRepository := c.MustGet(userRepositoryKey).(database.IUserRepository)
	messageRepository := c.MustGet(messageRepositoryKey).(database.IMessageRepository)

	username := c.Param("username")
	profileUser, err := userRepository.GetByUsername(username)
	if (profileUser == models.User{}) || err != nil {
		c.JSON(404, nil) // abort(404)
		return
	}

	followed := false
	u, userLoggedIn := c.Get("user")
	user := models.User{}

	if userLoggedIn { // if g.user from .py
		user = u.(models.User)
		followed, _ = userRepository.IsFollowing(user.UserId, profileUser.UserId)
	}

	messages, err := messageRepository.GetByUserId(profileUser.UserId, perPage)
	if err != nil {
		log.Fatal(err)
	}

	timelineData := TimelineData{
		IsPublicTimeline: false,
		IsMyTimeline:     user.Username == profileUser.Username,
		IsFollowed:       followed,
		HasMessages:      len(messages) > 0,

		ProfileUser: profileUser,

		Messages: messages,
	}

	renderTemplate(c, "timeline.html", &timelineData)
}

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

	c.Redirect(302, timeLineUrl)
}

// Registers a new message for the user.
func addMessage(c *gin.Context) {
	user, userLoggedIn := c.Get("user")

	if !userLoggedIn {
		c.JSON(401, nil)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	form := c.Request.Form
	text := form.Get("text")

	if text != "" {
		messageRepository := c.MustGet(messageRepositoryKey).(database.IMessageRepository)
		messageRepository.Create(models.Message{
			Author:  user.(models.User),
			Text:    text,
			PubDate: time.Now().Unix(),
		})

		flash(c, "Your message was recorded")
	}

	c.Redirect(302, timeLineUrl)
}

func loginGet(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(302, timeLineUrl)
		return
	}

	renderTemplate(c, "login.html", &LoginData{})
}

// Logs the user in.
func loginPost(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(302, timeLineUrl)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	form := c.Request.Form

	username, password := form.Get("username"), form.Get("password")
	userRepository := c.MustGet(userRepositoryKey).(database.IUserRepository)
	user, err := userRepository.GetByUsername(username)
	fmt.Println(user)
	if err != nil {
		log.Fatal(err)
	}

	session := sessions.Default(c)
	var errMsg string
	if (models.User{} == user) {
		errMsg = "Invalid username"
	} else if !pwdHash.CheckPasswordHash(password, user.PasswordHash) {
		errMsg = "Invalid password"
	} else {
		session.Set("user_id", user.UserId)
		session.Save()

		flash(c, "You were logged in")

		c.Redirect(302, timeLineUrl)
		return
	}

	renderTemplate(c, "login.html", &LoginData{
		Username: username,
		ErrorMsg: errMsg,
	})
}

// Shows the page for registering the user.
func registerGet(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(302, timeLineUrl)
		return
	}

	renderTemplate(c, "register.html", &RegisterData{})
}

// Registers the user.
func registerPost(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		// c.Writer.WriteHeader(http.StatusNoContent)
		c.Redirect(302, timeLineUrl)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	form := c.Request.Form

	userRepository := c.MustGet(userRepositoryKey).(database.IUserRepository)

	var errMsg string
	if form.Get("username") == "" {
		errMsg = "You have to enter a username"
	} else if form.Get("email") == "" || !strings.Contains(form.Get("email"), "@") {
		errMsg = "You have to enter a valid email address"
	} else if form.Get("password") == "" {
		errMsg = "You have to enter a password"
	} else if form.Get("password") != form.Get("password2") {
		errMsg = "The two passwords do not match"
	} else if user, _ := userRepository.GetByUsername(form.Get("username")); (user != models.User{}) { // TODO: This is heresy
		errMsg = "The username is already taken"
	} else {
		userRepository.Create(models.User{
			Username:     form.Get("username"),
			Email:        form.Get("email"),
			PasswordHash: pwdHash.GeneratePasswordHash(form.Get("password")),
		})
		if err != nil {
			log.Fatal(err)
		}

		flash(c, "You were successfully registered and can login now")

		c.Redirect(302, loginUrl)
		return
	}

	renderTemplate(c, "register.html", &RegisterData{
		ErrorMsg: errMsg,
		Username: form.Get("username"),
		Email:    form.Get("email"),
	})
}

// Logs the user out
func logout(c *gin.Context) {
	flash(c, "You were logged out")
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()
	c.Redirect(302, publicTimelineUrl)
}

func renderTemplate(c *gin.Context, templateSubPath string, templateData DataProvider) {
	templateData.initLayoutData()
	templateData.setFlashes(getFlashes(c))
	if user, userExists := c.Get("user"); userExists {
		templateData.setUser(user.(models.User))
	}

	path := templatePath(templateSubPath)
	t := parseHtmlFiles(path)
	err := t.Execute(c.Writer, templateData)
	if err != nil {
		log.Fatal(err)
	}
}

func parseHtmlFiles(files ...string) *template.Template {
	files = append(files, "./src/templates/layout.html")

	name := path.Base(files[0])
	t, err := template.New(name).Funcs(getHtmlDefaults()).ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func getHtmlDefaults() template.FuncMap {
	return template.FuncMap{
		"formatDatetime": formatDateTime,
		"gravatarUrl":    gravatarUrl,
	}
}

// Prepends the path to the templates folder to the given path.
func templatePath(subPath string) string {
	return fmt.Sprintf("%s/%s", "./src/templates", subPath)
}
