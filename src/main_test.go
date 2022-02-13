package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testDatabase = "./minitwit-test.db"
)

type TestClient struct {
	Router *gin.Engine
}

func init() {
	database = testDatabase
	os.Chdir("..")
}

func TestMain(m *testing.M) {
	initDb()
	defer os.Remove(testDatabase)
	exitVal := m.Run()
	os.Exit(exitVal)
}

// Helper function to register a user
func testUtilRegister(c *TestClient, username string, password string, password2 string, email string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("password2", password2)
	data.Set("email", email)
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))

	c.Router.ServeHTTP(w, req)

	return w
}

// Helper function to login
func testUtilLogin(c *TestClient, username string, password string) *httptest.ResponseRecorder {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))

	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	return w
}

// Registers and logs in in one go
func testUtilRegisterAndLogin(c *TestClient, username string, password string) *httptest.ResponseRecorder {
	testUtilRegister(c, username, password, "", "")
	return testUtilLogin(c, username, password)
}

// Helper function to logout
func testUtilLogout(c *TestClient) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/logout", nil)
	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	return w
}

// Records a message
func testUtilAddMessage(c *TestClient, text string) *httptest.ResponseRecorder {
	data := url.Values{}
	data.Set("text", text)
	req, _ := http.NewRequest("POST", "/add_message", strings.NewReader(data.Encode()))
	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	return w
}

// Make sure registering works
func TestRegister(t *testing.T) {
	c := &TestClient{Router: setupRouter()}

	rv := testUtilRegister(c, "user1", "default", "", "")
	assert.Contains(t, rv.Body.String(), "You were successfully registered and can login now")

	rv = testUtilRegister(c, "user1", "default", "", "")
	assert.Contains(t, rv.Body.String(), "The username is already taken")

	rv = testUtilRegister(c, "", "default", "", "")
	assert.Contains(t, rv.Body.String(), "You have to enter a username")

	rv = testUtilRegister(c, "meh", "", "", "")
	assert.Contains(t, rv.Body.String(), "The have to enter a password")

	rv = testUtilRegister(c, "meh", "x", "y", "")
	assert.Contains(t, rv.Body.String(), "The two passwords do not match")

	rv = testUtilRegister(c, "meh", "foo", "", "broken")
	assert.Contains(t, rv.Body.String(), "You have to enter a valid email address")
}

// Make sure logging in and logging out works
func TestLoginLogout(t *testing.T) {
	c := &TestClient{Router: setupRouter()}
	rv := testUtilRegisterAndLogin(c, "user1", "default")
	assert.Contains(t, rv.Body.String(), "You were logged in")
	rv = testUtilLogout(c)
	assert.Contains(t, rv.Body.String(), "You were logged out")
	rv = testUtilLogin(c, "user1", "wrongpassword")
	assert.Contains(t, rv.Body.String(), "Invalid password")
	rv = testUtilLogin(c, "user2", "wrongpassword")
	assert.Contains(t, rv.Body.String(), "Invalid username")
}

// Check if adding messages works
func TestMessageRecording(t *testing.T) {
	c := &TestClient{Router: setupRouter()}

	testUtilRegisterAndLogin(c, "foo", "default")
	testUtilAddMessage(c, "test message 1")
	testUtilAddMessage(c, "<test message 2>")
	req, _ := http.NewRequest("get", "/", nil)
	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "test message 1")
	assert.Contains(t, w.Body.String(), "&lt;test message 2&gt;")
}

// Make sure that timelines work
func TestTimelines(t *testing.T) {
	c := &TestClient{Router: setupRouter()}

	testUtilRegisterAndLogin(c, "foo", "default")
	testUtilAddMessage(c, "the message by foo")
	testUtilLogout(c)
	testUtilRegisterAndLogin(c, "bar", "default")
	testUtilAddMessage(c, "the message by bar")
	rv, _ := http.NewRequest("GET", "/public", nil)
	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.Contains(t, w.Body.String(), "the message by foo")
	assert.Contains(t, w.Body.String(), "the message by bar")

	// bar's timeline should just show bar's message
	rv, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.NotContains(t, w.Body.String(), "the message by foo")
	assert.Contains(t, w.Body.String(), "the message by bar")

	// now let's follow foo
	rv, _ = http.NewRequest("GET", "/foo/follow", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.Contains(t, w.Body.String(), "You are now following &#34;foo&#34;")

	// we should now see foo's message
	rv, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.Contains(t, w.Body.String(), "the message by foo")
	assert.Contains(t, w.Body.String(), "the message by bar")

	// but on the user's page we only want the user's message
	rv, _ = http.NewRequest("GET", "/bar", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.NotContains(t, w.Body.String(), "the message by foo")
	assert.Contains(t, w.Body.String(), "the message by bar")
	rv, _ = http.NewRequest("GET", "/foo", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.Contains(t, w.Body.String(), "the message by foo")
	assert.NotContains(t, w.Body.String(), "the message by bar")

	// now unfollow and check if that worked
	rv, _ = http.NewRequest("GET", "/foo/unfollow", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.Contains(t, w.Body.String(), "You are no longer following &#34;foo&#34;")

	rv, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	c.Router.ServeHTTP(w, rv)
	assert.NotContains(t, w.Body.String(), "the message by foo")
	assert.Contains(t, w.Body.String(), "the message by bar")
}
