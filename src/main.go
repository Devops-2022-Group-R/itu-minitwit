package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
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
)

const (
	perPage   = 30
	debug     = true
	secretKey = "development key"

	timeLineUrl       = "/"
	publicTimelineUrl = "/public"
	loginUrl          = "/login"
	flashesKey = "flashes"
)

var databasePath = "./minitwit.db"

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
	r.POST("/add_message", addMessage)
	r.GET(loginUrl, loginGet)
	r.POST(loginUrl, controllers.LoginPost)
	r.POST("/register", controllers.RegisterController)
	r.GET("/logout", logout)

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
	db := connectDb()
	defer db.Close()

	file, err := ioutil.ReadFile("./src/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	query := string(file)

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// Convenience method to look up the id for a username.
func getUserId(username string, db *sql.DB) *int64 {
	row := db.QueryRow("select user_id from user where username = ?", username)

	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		return nil
	}

	return &userId
}

func getUserFromUsername(username string, db *sql.DB) *User {
	rows := database.QueryDb(db, "select * from user where username = ?", username)

	return parseUser(rows)
}

func getUserFromId(id int64, db *sql.DB) *User {
	rows := database.QueryDb(db, "select * from user where user_id = ?", id)

	return parseUser(rows)
}

func parseUser(rows []map[string]interface{}) *User {
	if len(rows) == 0 {
		return nil
	}

	return &User{
		UserId:       rows[0]["user_id"].(int64),
		Username:     rows[0]["username"].(string),
		Email:        rows[0]["email"].(string),
		PasswordHash: rows[0]["pw_hash"].(string),
	}
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
	db := connectDb()
	defer db.Close()

	c.Set("db", db)

	session := sessions.Default(c)
	if userId := session.Get("user_id"); userId != nil {
		user := *getUserFromId(userId.(int64), db)
		c.Set("user", user)
	}

	c.Next()
}

// Shows a users timeline or if no user is logged in it will redirect to the
// public timeline. This timeline shows the user's messages as well as all the
// messages of followed users.
func timeline(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	defer db.Close()

	var user User
	if u, isLoggedIn := c.Get("user"); isLoggedIn {
		user = u.(User)
	} else {
		c.Redirect(302, publicTimelineUrl)
		return
	}

	// c.Request.args.Get("offset", int) // offset = request.args.get('offset', type=int)
	query :=
		`
			select message.*, user.* from message, user
			where message.flagged = 0 and message.author_id = user.user_id and 
			(
				user.user_id = ? or
				user.user_id in 
				(
					select whom_id from follower
					where who_id = ?
				)
			)
			order by message.pub_date desc limit ?
		`
	results := database.QueryDb(db, query, user.UserId, user.UserId, perPage) //  [session['user_id'], session['user_id'], PER_PAGE]))

	messages := createTweetsFromQuery(results)

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
	db := c.MustGet("db").(*sql.DB)
	defer db.Close()

	query := `
		select message.*, user.* 
		from message, user 
		where message.flagged = 0 
			and message.author_id = user.user_id 
		order by message.pub_date desc
		limit ?`
	results := database.QueryDb(db, query, perPage)

	messages := make([]Message, 0)

	for _, result := range results {
		message := Message{
			Email:    result["email"].(string),
			Username: result["username"].(string),
			Text:     result["text"].(string),
			PubDate:  result["pub_date"].(int64),
		}

		messages = append(messages, message)
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
	db := c.MustGet("db").(*sql.DB)
	defer db.Close()

	username := c.Param("username")

	profileUser := getUserFromUsername(username, db)
	if profileUser == nil {
		c.JSON(404, nil) // abort(404)
		return
	}

	followed := false
	u, userLoggedIn := c.Get("user")
	user := User{
		Username: "",
	}

	if userLoggedIn { // if g.user from .py
		user = u.(User)
		query :=
			`
				select 1 
				from follower 
				where follower.who_id = ? and follower.whom_id = ?
			`
		session := sessions.Default(c)
		userId := session.Get("user_id").(int64)
		isFollowing := database.QueryDb(db, query, userId, profileUser.UserId) // [session['user_id'], profile_user['user_id']], one=True) is not None
		if len(isFollowing) > 0 {                                     // this condition is trying to check -> "is not none" - is this correct?
			followed = true
		}
	}

	messageQuery :=
		`
			select message.*, user.* from message, 
			user where user.user_id = message.author_id and user.user_id = ?
			order by message.pub_date desc limit ?
		`
	results := database.QueryDb(db, messageQuery, profileUser.UserId, perPage)
	messages := createTweetsFromQuery(results)

	timelineData := TimelineData{
		IsPublicTimeline: false,
		IsMyTimeline:     user.Username == profileUser.Username,
		IsFollowed:       followed,
		HasMessages:      len(messages) > 0,

		ProfileUser: *profileUser,

		Messages: messages,
	}

	renderTemplate(c, "timeline.html", &timelineData)
}

// Convenience transforming data from a query into messages
func createTweetsFromQuery(results []map[string]interface{}) []Message {
	messages := make([]Message, 0)

	for _, result := range results {
		message := Message{
			Email:    result["email"].(string),
			Username: result["username"].(string),
			PubDate:  result["pub_date"].(int64),
			Text:     result["text"].(string),
		}

		messages = append(messages, message)
	}
	return messages
}

// Adds the current user as follower of the given user.
func followUser(c *gin.Context) {
	username := c.Param("username")
	db := c.MustGet("db").(*sql.DB)
	whomID := database.GetUserId(username, db)
	session := sessions.Default(c)

	if _, isLoggedIn := c.Get("user"); !isLoggedIn {
		c.JSON(401, nil)
		return
	}
	if whomID == nil {
		c.JSON(404, nil)
		return
	}

	database.QueryDb(db, "insert into follower (who_id, whom_id) values (?, ?)",
		session.Get("user_id"), whomID)

	flash(c, fmt.Sprintf("You are now following %s", username))

	c.Redirect(302, timeLineUrl)
}

// Removes the current user as follower of the given user.
func unfollowUser(c *gin.Context) {
	username := c.Param("username")
	db := c.MustGet("db").(*sql.DB)
	whomID := database.GetUserId(username, db)
	session := sessions.Default(c)

	if _, isLoggedIn := c.Get("user"); !isLoggedIn {
		c.JSON(401, nil)
		return
	}
	if whomID == nil {
		c.JSON(404, nil)
		return
	}

	database.QueryDb(db, "delete from follower where who_id = ? and whom_id = ?",
		session.Get("user_id"), whomID)

	flash(c, fmt.Sprintf("You are no longer following %s", username))

	c.Redirect(302, timeLineUrl)
}

// Registers a new message for the user.
func addMessage(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

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
		database.QueryDb(db, "insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, 0)",
			user.(User).UserId, text, time.Now().UTC().Unix())

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
		templateData.setUser(user.(User))
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
