package controllers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/controllers"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
)

type BaseTestSuite struct {
	suite.Suite
	openDatabase func() gorm.Dialector
}

func (suite *BaseTestSuite) SetupTest() {
	db := sqlite.Open("file::memory:?cache=shared")
	suite.openDatabase = func() gorm.Dialector {
		return db
	}
	database.InitDatabase(suite.openDatabase)
}

func (suite *BaseTestSuite) TearDownTest() {
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
	suite.Run(t, new(BaseTestSuite))
}

// Helpers
func (suite *BaseTestSuite) sendRequest(req *http.Request) *httptest.ResponseRecorder {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	router := controllers.SetupRouter(suite.openDatabase)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func (suite *BaseTestSuite) sendAuthedRequest(req *http.Request, username string, password string) *httptest.ResponseRecorder {
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64(username, password))

	return suite.sendRequest(req)
}

func (suite *BaseTestSuite) sendSimulatorRequest(req *http.Request) *httptest.ResponseRecorder {
	return suite.sendAuthedRequest(req, "simulator", "super_safe!")
}

func (suite *BaseTestSuite) registerUser(username, email, password string) {
	body, _ := json.Marshal(gin.H{"username": username, "email": email, "pwd": password})
	suite.sendRequest(httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body)))
}

func (suite *BaseTestSuite) registerSimulator() {
	suite.registerUser("simulator", "simulator@simulator.dk", "super_safe!")
}

func (suite *BaseTestSuite) readBody(w *httptest.ResponseRecorder, target interface{}) {
	bodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(bodyBytes, target)
}

func encodeCredentialsToB64(username string, password string) string {
	data := username + ":" + password
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}
