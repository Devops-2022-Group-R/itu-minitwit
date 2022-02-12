package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	database  = "./minitwit.db"
	perPage   = 30
	debug     = true
	secretKey = "development key"

	timeLineUrl       = "/"
	publicTimelineUrl = "/public"
	loginUrl          = "/login"
)

type Row = map[string]interface{}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Use(beforeRequest)

	store := cookie.NewStore([]byte(secretKey))
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/static", "static")

	r.GET(loginUrl, loginGetHandler)
	r.POST(loginUrl, loginPostHandler)
	r.GET(publicTimelineUrl, publicTimeline)
	r.POST("/add_message", addMessage)
	r.GET("/:username", userTimeline)

	r.Run()
}

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./minitwit.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Creates the database tables.
func initDb() {
	db := connectDb()
	defer db.Close()

	file, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	query := string(file)

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// Queries the database and returns a list of dictionaries.
func queryDb(db *sql.DB, query string, args ...interface{}) []Row {
	rows, err := db.Query(query, args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	results := make([]Row, 0)
	for rows.Next() {
		row := make(Row)
		for i := 0; i < len(columnNames); i++ {
			key := columnNames[i]
			var value interface{}
			err := rows.Scan(&value)
			if err != nil {
				log.Fatal(err)
			}
			row[key] = value
		}
		results = append(results, row)
	}
	err = rows.Err()

	return results
}

// Convenience method to look up the id for a username.
func getUserId(username string, db *sql.DB) *int {
	row := db.QueryRow("select user_id from user where username = ?", username)

	var userId int
	err := row.Scan(&userId)
	if err != nil {
		return nil
	}

	return &userId
}

// Format a timestamp for display.
func formatDateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("%Y-%m-%d @ %H:%M")
}

// Return the gravatar image for the given email address.
func gravatarUrl(email string, size int) string {
	email = strings.ToLower(strings.TrimSpace(email))

	hash := md5.New()
	hash.Write([]byte(email))

	hex := fmt.Sprintf("%x", hash.Sum(nil))

	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hex, size)
}

// Make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(c *gin.Context) {
	db := connectDb()
	defer db.Close()

	c.Set("db", db)
	c.Set("user", nil)

	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		users := queryDb(db, "select * from user where user_id = ?", session.Get("user_id"))
		c.Set("user", users[0])
	}

	c.Next()
}

// @app.route('/')
func timeline(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	session := sessions.Default(c)
	defer db.Close()

	if _, isLoggedIn := c.Get("user"); !isLoggedIn {
		c.Redirect(307, publicTimelineUrl)
	}

	userId := session.Get("user_id").(int)
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
	results := queryDb(db, query, userId, userId, perPage) //  [session['user_id'], session['user_id'], PER_PAGE]))

	messages := make([]Message, 0)
	users := make([]User, 0)

	for _, result := range results {
		message := Message{
			Email:    result["email"].(string),
			Username: result["username"].(string),
			PubDate:  result["pub_date"].(int),
			Text:     result["text"].(string),
		}

		user := User{
			Username: result["username"].(string),
		}

		messages = append(messages, message)
		users = append(users, user)
	}

	renderTemplate(c, "timeline.html", TimelineData{
		IsPublicTimeline: false,
		IsMyTimeline:     false,
		IsFollowed:       false,
		HasMessages:      len(messages) > 0,
		Messages:         messages,
	})
}

// Displays the latest messages of all users.
// @app.route('/public')
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
	results := queryDb(db, query, perPage)

	messages := make([]Message, 0)
	users := make([]User, 0)

	for _, result := range results {
		message := Message{
			Email:    result["email"].(string),
			Username: result["username"].(string),
			Text:     result["text"].(string),
			PubDate:  result["pub_date"].(int),
		}
		user := User{
			Username: result["username"].(string),
		}

		messages = append(messages, message)
		users = append(users, user)
	}

	timelineData := TimelineData{
		IsPublicTimeline: true,
		IsMyTimeline:     false,
		IsFollowed:       true,
		HasMessages:      len(messages) > 0,
		Messages:         messages,
	}
}

// Display's a users tweets.
func userTimeline(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	defer db.Close()

	username := c.Param("username")
	queryUserProfile := "select * from user where username = ?"
	profileUser := queryDb(db, queryUserProfile, username)

	if _, contain := profileUser[0][username]; !contain {
		c.JSON(404, nil) // abort(404)
	}
	followed := false

	if _, userExists := c.Get("user"); userExists{ // if g.user from .py
		query :=
			`
				"select 1 
				from follower 
				where follower.who_id = ? and follower.whom_id = ?"
				`
		session := sessions.Default(c)
		userId := session.Get("user_id").(int)
		isFollowing := queryDb(db, query, userId, profileUser[0]["user_id"]) // [session['user_id'], profile_user['user_id']], one=True) is not None
		if len(isFollowing[0]) > 0 {                                         // this condition is trying to check -> "is not none" - is this correct?
			followed = true
		}
	}

	messageQuery :=
		`
			select message.*, user.* from message, 
			user where user.user_id = message.author_id and user.user_id = ?
			order by message.pub_date desc limit
		`
	results := queryDb(db, messageQuery, profileUser[0]["user_id"], perPage)

	messages := make([]Message, 0)
	users := make([]User, 0)


	timelineData := TimelineData{
		IsPublicTimeline: false,
		IsMyTimeline: ,
		IsFollowed:       followed,
		HasMessages: len(messages) > 0,
		User: ,
		Messages: messages,
	}

	renderTemplate(c, "timeline.html",timelineData)
}

// Adds the current user as follower of the given user.
// @app.route('/<username>/follow')
func followUser(username string, c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	user := c.MustGet("user").(Row)
	whomID := getUserId(username, c)
	session := sessions.Default(c)

	if user == nil {
		c.JSON(401, nil)
	}
	if whomID == nil {
		c.JSON(404, nil)
	}

	queryDb(db, "insert into follower (who_id, whom_id) values (?, ?)",
		session.Get("user_id"), whomID)

	// TODO: flash('You are now following "%s"' % username)

	c.Redirect(307, timeLineUrl)
}

// Removes the current user as follower of the given user.
// @app.route('/<username>/unfollow')
func unfollowUser(username string, c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	user := c.MustGet("user").(Row)
	whomID := getUserId(username, c)
	session := sessions.Default(c)

	if user == nil {
		c.JSON(401, nil)
	}
	if whomID == nil {
		c.JSON(404, nil)
	}

	queryDb(db, "delete from follower where who_id = ? and whom_id = ?",
		session.Get("user_id"), whomID)

	// TODO: flash('You are no longer following "%s"' % username)

	c.Redirect(307, timeLineUrl)
}

// Registers a new message for the user.
func addMessage(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(401, nil)
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	form := c.Request.Form
	text := form.Get("text")

	if text != "" {
		queryDb(db, "insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, 0)",
			session.Get("user_id"), text, time.Now())

		// TODO: flash('Your message was recorded')
	}

	// TODO: return redirect(url_for('timeline'))
}

func loginGetHandler(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(307, timeLineUrl)
		return
	}

	renderTemplate(c, "login.html", LoginData{})
}

// Logs the user in.
func loginPostHandler(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(307, timeLineUrl)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	form := c.Request.Form

	username, password := form.Get("username"), form.Get("password")
	db := c.MustGet("db").(*sql.DB)
	users := queryDb(db, "select * from user where username = ?", username)

	session := sessions.Default(c)
	var errMsg string
	if len(users) == 0 {
		errMsg = "Invalid username"
	} else if !checkPasswordHash(users[0]["pw_hash"], password) {
		errMsg = "Invalid password"
	} else {
		session.Set("user_id", users[0]["user_id"])
		// TODO: Translate this from Python - flash('You were logged in')
		c.Redirect(307, timeLineUrl)
		return
	}

	renderTemplate(c, "login.html", LoginData{
		Username: username,
		ErrorMsg: errMsg,
	})
}

// Shows the page for registering the user.
func registerGetHandler(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(307, timeLineUrl)
		return
	}

	renderTemplate(c, "register.html", RegisterData{})
}

// Registers the user.
func registerPostHandler(c *gin.Context) {
	if _, userIsInSession := c.Get("user"); userIsInSession {
		c.Redirect(307, timeLineUrl)
		return
	}

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	form := c.Request.Form

	db := c.MustGet("db").(*sql.DB)
	var errMsg string
	if form.Get("username") == "" {
		errMsg = "You have to enter a username"
	} else if form.Get("email") == "" || !strings.Contains(form.Get("email"), "@") {
		errMsg = "You have to enter a valid email address"
	} else if form.Get("pasword") == "" {
		errMsg = "You have to enter a password"
	} else if form.Get("pasword") != form.Get("pasword2") {
		errMsg = "The two passwords do not match"
	} else if getUserId(form.Get("username"), db) != nil {
		errMsg = "The username is already taken"
	} else {
		_, err := db.Exec(
			"insert into user (username, email, pw_hash) values (?, ?, ?)",
			form.Get("username"),
			form.Get("email"),
			generatePasswordHash(form.Get("password")),
			// TODO: Translate from Python - flash('You were successfully registered and can login now')
		)
		if err != nil {
			log.Fatal(err)
		}

		c.Redirect(307, loginUrl)
	}

	renderTemplate(c, "register.html", RegisterData{
		ErrorMsg: errMsg,
		Username: form.Get("username"),
		Email:    form.Get("email"),
	})
}

// Logs the user out
func logout(c *gin.Context) {
	// TODO: Translate this from Python - flash('You were logged out')
	session := sessions.Default(c)
	session.Delete("user_id")
	c.Redirect(307, publicTimelineUrl)
}

func renderTemplate(c *gin.Context, templateSubPath string, templateData interface{}) {
	path := templatePath(templateSubPath)
	t := parseHtmlFiles(path)
	err := t.Execute(c.Writer, templateData)
	if err != nil {
		log.Fatal(err)
	}
}

func parseHtmlFiles(files ...string) *template.Template {
	files = append(files, "./new-templates/layout.html")

	name := path.Base(files[0])
	t, err := template.New(name).Funcs(getHtmlDefaults()).ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func getHtmlDefaults() template.FuncMap {
	return template.FuncMap{
		// "getStaticRoute": getStaticRoute,
		"formatDatetime": formatDateTime,
		"gravatarUrl":    gravatarUrl,
	}
}

// Prepends the path to the templates folder to the given path.
func templatePath(subPath string) string {
	return fmt.Sprintf("%s/%s", "./src/templates", subPath)
}
