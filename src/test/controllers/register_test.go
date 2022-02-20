package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct{ BaseTestSuite }

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (suite *RegisterTestSuite) TestRegisterController_GivenNoBody_Returns400() {
	req, _ := http.NewRequest(http.MethodPost, "/register", nil)
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenMissingField_Returns400() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"pwd":      "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenValidRequest_Returns204() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yennefer@aretuza.wr",
		"pwd":      "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenInvalidEmail_Returns422() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yenneferaretuza.wr",
		"pwd":      "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_RunTwiceWithSameUsername_Returns409() {
	assert := assert.New(suite.T())

	firstRegister, _ := json.Marshal(gin.H{
		"username": "GeraltLover",
		"email":    "yennefer@aretuza.wr",
		"pwd":      "chaosmaster",
	})
	secondRegister, _ := json.Marshal(gin.H{
		"username": "GeraltLover",
		"email":    "triss@merigold.wr",
		"pwd":      "peacemaster",
	})

	req1, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(firstRegister))
	w1 := suite.sendRequest(req1)

	assert.Equal(http.StatusNoContent, w1.Code)

	req2, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(secondRegister))
	w2 := suite.sendRequest(req2)

	assert.Equal(http.StatusConflict, w2.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenValidBody_AddsUserToDatabase() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yennefer@aretuza.wr",
		"pwd":      "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	suite.sendRequest(req)

	gormDb, _ := database.ConnectDatabase(suite.openDatabase)
	userRepo := database.NewGormUserRepository(gormDb)
	user, err := userRepo.GetByUsername("Yennefer of Vengerberg")

	assert := assert.New(suite.T())
	assert.Nil(err)
	assert.Equal("Yennefer of Vengerberg", user.Username)
	assert.Equal("yennefer@aretuza.wr", user.Email)
	assert.NotEmpty(user.PasswordHash)
}
