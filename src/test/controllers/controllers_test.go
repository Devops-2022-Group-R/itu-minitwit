package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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

func (suite *TestSuite) sendRequest(req *http.Request) *httptest.ResponseRecorder {
	router := controllers.SetupRouter(suite.openDatabase)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestControllersTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
