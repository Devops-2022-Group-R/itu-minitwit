package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
)

const followUrl = "/fllws/geralt"

// Unsafe to use unless you are strictly using the username
func utilCreateUsersInDatabase(suite *TestSuite) {
	gormDb, _ := database.ConnectDatabase(suite.openDatabase)
	userRepo := database.NewGormUserRepository(gormDb)

	userRepo.Create(models.User{UserId: 1, Username: "geralt", Email: "geralt@witcher.pl", PasswordHash: "123"})
	userRepo.Create(models.User{UserId: 2, Username: "yennefer", Email: "yennefer@witcher.pl", PasswordHash: "123"})
	userRepo.Create(models.User{UserId: 3, Username: "triss", Email: "triss@witcher.pl", PasswordHash: "123"})
}

func (suite *TestSuite) TestFollowPostController_GivenValidBody_Returns204() {
	utilCreateUsersInDatabase(suite)

	body, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenNonExistingUser_Returns404() {
	body, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req, _ := http.NewRequest(http.MethodPost, "/fllws/i-dont-exist", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenFollowANonExistingUser_Returns404() {
	utilCreateUsersInDatabase(suite)

	body, _ := json.Marshal(gin.H{"follow": "vesemir"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestFollowGetController_GivenUserWithNoFollows_ReturnsEmpty() {
	utilCreateUsersInDatabase(suite)

	reqBody, _ := json.Marshal(gin.H{"follow": "vesemir"})
	req, _ := http.NewRequest(http.MethodGet, followUrl, bytes.NewReader(reqBody))
	w := suite.sendRequest(req)

	var resBody []string
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Empty(suite.T(), resBody)
}

func (suite *TestSuite) TestFollowGetController_GivenUserWithFollows_ReturnsAllFollowed() {
	utilCreateUsersInDatabase(suite)
	assert := assert.New(suite.T())

	// Create follow relationships
	body1, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req1, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body1))
	w1 := suite.sendRequest(req1)
	body2, _ := json.Marshal(gin.H{"follow": "triss"})
	req2, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body2))
	w2 := suite.sendRequest(req2)
	assert.Equal(http.StatusNoContent, w1.Code)
	assert.Equal(http.StatusNoContent, w2.Code)

	// Get users
	reqBody, _ := json.Marshal(gin.H{"follow": "vesemir"})
	req, _ := http.NewRequest(http.MethodGet, followUrl, bytes.NewReader(reqBody))
	w := suite.sendRequest(req)

	var resBody []string
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Equal(http.StatusOK, w.Code)
	assert.ElementsMatch([...]string{"yennefer", "triss"}, resBody)
}
