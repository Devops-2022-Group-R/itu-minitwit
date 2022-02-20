package controllers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
)

type TestSuite struct {
	suite.Suite
	openDatabase func() gorm.Dialector
}

func (suite *TestSuite) SetupTest() {
	db := sqlite.Open("file::memory:?cache=shared")
	suite.openDatabase = func() gorm.Dialector {
		return db
	}
	database.InitDatabase(suite.openDatabase)
}

func (suite *TestSuite) TearDownTest() {
	// Clear everything from the database
	db, _ := database.ConnectDatabase(suite.openDatabase)
	db.Exec(`
		PRAGMA writable_schema = 1;
		delete from sqlite_master where type in ('table', 'index', 'trigger');
		PRAGMA writable_schema = 0;
		VACUUM;
	`)
}

func TestControllersTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// Helpers
func (suite *TestSuite) sendRequest(req *http.Request) *httptest.ResponseRecorder {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	router := controllers.SetupRouter(suite.openDatabase)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func (suite *TestSuite) sendAuthedRequest(req *http.Request, username string, password string) *httptest.ResponseRecorder {
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64(username, password))

	return suite.sendRequest(req)
}

func (suite *TestSuite) sendSimulatorRequest(req *http.Request) *httptest.ResponseRecorder {
	return suite.sendAuthedRequest(req, "simulator", "super_safe!")
}

func (suite *TestSuite) registerUser(username, email, password string) {
	body, _ := json.Marshal(gin.H{"username": username, "email": email, "pwd": password})
	suite.sendRequest(httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body)))
}

func (suite *TestSuite) registerSimulator() {
	suite.registerUser("simulator", "simulator@simulator.dk", "super_safe!")
}

func encodeCredentialsToB64(username string, password string) string {
	data := username + ":" + password
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}
