package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDatabase = "./minitwit-test.db"
)

func init() {
	os.Chdir("..")
}

func TestMain(m *testing.M) {
	initDb()
	exitVal := m.Run()
	os.Remove(testDatabase)
	os.Exit(exitVal)
}

func testUtilNewHttpClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	return client
}

func testUtilGetBody(body io.ReadCloser) string {
	b, _ := ioutil.ReadAll(body)
	return string(b)
}

type TestClient struct {
	Server *httptest.Server
	Client *http.Client
}

func (c *TestClient) Close() {
	c.Server.Close()
}

func (c *TestClient) Url(path string) string {
	return fmt.Sprintf("%s%s", c.Server.URL, path)
}

func testUtilNewTestClient() *TestClient {
	return &TestClient{
		Server: httptest.NewServer(setupRouter()),
		Client: testUtilNewHttpClient(),
	}
}

// Helper function to register a user
func testUtilRegister(client *TestClient, username string, password string, password2 string, email string) *http.Response {
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
	req, _ := http.NewRequest(http.MethodPost, client.Url("/register"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Client.Do(req)
	return resp
}

// Helper function to login
func testUtilLogin(c *TestClient, username string, password string) *http.Response {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	req, _ := http.NewRequest("POST", c.Url("/login"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := c.Client.Do(req)
	return resp
}

// Registers and logs in in one go
func testUtilRegisterAndLogin(c *TestClient, username string, password string) *http.Response {
	testUtilRegister(c, username, password, "", "")
	return testUtilLogin(c, username, password)
}

// Helper function to logout
func testUtilLogout(c *TestClient) *http.Response {
	req, _ := http.NewRequest("GET", c.Url("/logout"), nil)
	resp, _ := c.Client.Do(req)
	return resp
}

// Records a message
func testUtilAddMessage(c *TestClient, text string) *http.Response {
	data := url.Values{}
	data.Set("text", text)
	req, _ := http.NewRequest("POST", c.Url("/add_message"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := c.Client.Do(req)
	return resp
}

// Make sure registering works
func TestRegister(t *testing.T) {
	ts := testUtilNewTestClient()
	defer ts.Close()

	rv := testUtilRegister(ts, "user1", "default", "", "")
	assert.Contains(t, testUtilGetBody(rv.Body), "You were successfully registered and can login now")

	rv = testUtilRegister(ts, "user1", "default", "", "")
	assert.Contains(t, testUtilGetBody(rv.Body), "The username is already taken")

	rv = testUtilRegister(ts, "", "default", "", "")
	assert.Contains(t, testUtilGetBody(rv.Body), "You have to enter a username")

	rv = testUtilRegister(ts, "meh", "", "", "")
	assert.Contains(t, testUtilGetBody(rv.Body), "You have to enter a password")

	rv = testUtilRegister(ts, "meh", "x", "y", "")
	assert.Contains(t, testUtilGetBody(rv.Body), "The two passwords do not match")

	rv = testUtilRegister(ts, "meh", "foo", "", "broken")
	assert.Contains(t, testUtilGetBody(rv.Body), "You have to enter a valid email address")
}

// Make sure logging in and logging out works
func TestLoginLogout(t *testing.T) {
	ts := testUtilNewTestClient()
	defer ts.Close()

	rv := testUtilRegisterAndLogin(ts, "user1", "default")
	assert.Contains(t, testUtilGetBody(rv.Body), "You were logged in")

	rv = testUtilLogout(ts)
	assert.Contains(t, testUtilGetBody(rv.Body), "You were logged out")

	rv = testUtilLogin(ts, "user1", "wrongpassword")
	assert.Contains(t, testUtilGetBody(rv.Body), "Invalid password")

	rv = testUtilLogin(ts, "user2", "wrongpassword")
	assert.Contains(t, testUtilGetBody(rv.Body), "Invalid username")
}

// Check if adding messages works
func TestMessageRecording(t *testing.T) {
	ts := testUtilNewTestClient()
	defer ts.Close()

	testUtilRegisterAndLogin(ts, "foo", "default")
	testUtilAddMessage(ts, "test message 1")
	testUtilAddMessage(ts, "<test message 2>")
	req, _ := http.NewRequest(http.MethodGet, ts.Url("/"), nil)
	resp, _ := ts.Client.Do(req)

	body := testUtilGetBody(resp.Body)
	assert.Contains(t, body, "test message 1")
	assert.Contains(t, body, "&lt;test message 2&gt;")
}

// Make sure that timelines work
func TestTimelines(t *testing.T) {
	ts := testUtilNewTestClient()
	defer ts.Close()

	testUtilRegisterAndLogin(ts, "foo", "default")
	testUtilAddMessage(ts, "the message by foo")
	testUtilLogout(ts)
	testUtilRegisterAndLogin(ts, "bar", "default")
	testUtilAddMessage(ts, "the message by bar")
	req, _ := http.NewRequest("GET", ts.Url("/public"), nil)
	rv, _ := ts.Client.Do(req)
	body := testUtilGetBody(rv.Body)
	assert.Contains(t, body, "the message by foo")
	assert.Contains(t, body, "the message by bar")

	// bar's timeline should just show bar's message
	req, _ = http.NewRequest("GET", ts.Url("/"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.NotContains(t, body, "the message by foo")
	assert.Contains(t, body, "the message by bar")

	// now let's follow foo
	req, _ = http.NewRequest("GET", ts.Url("/foo/follow"), nil)
	rv, _ = ts.Client.Do(req)
	assert.Contains(t, testUtilGetBody(rv.Body), "You are now following foo")

	// we should now see foo's message
	req, _ = http.NewRequest("GET", ts.Url("/"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.Contains(t, body, "the message by foo")
	assert.Contains(t, body, "the message by bar")

	// but on the user's page we only want the user's message
	req, _ = http.NewRequest("GET", ts.Url("/bar"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.NotContains(t, body, "the message by foo")
	assert.Contains(t, body, "the message by bar")

	req, _ = http.NewRequest("GET", ts.Url("/foo"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.Contains(t, body, "the message by foo")
	assert.NotContains(t, body, "the message by bar")

	// now unfollow and check if that worked
	req, _ = http.NewRequest("GET", ts.Url("/foo/unfollow"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.Contains(t, body, "You are no longer following foo")

	req, _ = http.NewRequest("GET", ts.Url("/"), nil)
	rv, _ = ts.Client.Do(req)
	body = testUtilGetBody(rv.Body)
	assert.NotContains(t, body, "the message by foo")
	assert.Contains(t, body, "the message by bar")
}
