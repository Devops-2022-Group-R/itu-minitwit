package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct{ BaseTestSuite }

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (suite *RegisterTestSuite) TestRegisterController_GivenNoBody_Returns400() {
	// Act
	req, _ := http.NewRequest(http.MethodPost, "/register", nil)
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenMissingField_Returns400() {
	// Act
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"pwd":      "chaosmaster",
	})
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenValidRequest_Returns204() {
	// Act
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yennefer@aretuza.wr",
		"pwd":      "chaosmaster",
	})
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenInvalidEmail_Returns422() {
	// Act
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yenneferaretuza.wr",
		"pwd":      "chaosmaster",
	})
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_RunTwiceWithSameUsername_Returns409() {
	// Act
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
	req2, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(secondRegister))
	w2 := suite.sendRequest(req2)

	suite.Equal(http.StatusNoContent, w1.Code)
	suite.Equal(http.StatusConflict, w2.Code)
}

func (suite *RegisterTestSuite) TestRegisterController_GivenValidBody_AddsUserToDatabase() {
	// Act
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yennefer@aretuza.wr",
		"pwd":      "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	suite.sendRequest(req)

	// Assert
	gormDb, _ := database.ConnectDatabase(suite.openDatabase)
	userRepo := database.NewGormUserRepository(gormDb)
	user, err := userRepo.GetByUsername("Yennefer of Vengerberg")

	suite.Nil(err)
	suite.Equal("Yennefer of Vengerberg", user.Username)
	suite.Equal("yennefer@aretuza.wr", user.Email)
	suite.NotEmpty(user.PasswordHash)
}
