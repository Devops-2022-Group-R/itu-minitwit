package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
)

func (suite *TestSuite) Test_GetMessages_Returns_OK() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) Test_GetUserMessages_Returns_OK() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")

	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) Test_GetUserMessages_Given_NonExistent_User_Returns_NotFound() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) Test_PostUserMessage_Returns_NoContent() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")
	body, _ := json.Marshal(gin.H{"content": "Omnis vir lupus."})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Darrow", "Reaper"))
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) Test_PostUserMessage_Given_NonExistent_User_Returns_NotFound() {
	// Arrange
	body, _ := json.Marshal(gin.H{"content": "Omnis vir lupus."})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Darrow", "Reaper"))
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) Test_PostUserMessage_Given_Empty_Message_Returns_BadRequest() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")
	body, _ := json.Marshal(gin.H{"content": ""})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Darrow", "Reaper"))
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *TestSuite) Test_PostUserMessage_Given_Simulator_Returns_No_Content() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")
	suite.registerUser("simulator", "simulator@andromedus.com", "super_safe!")
	body, _ := json.Marshal(gin.H{"content": ""})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("simulator", "super_safe!"))
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *TestSuite) Test_GetFeedMessages_Given_NonExistent_User_Returns_NotFound() {
	// Arrange
	encodedCredentials := encodeCredentialsToB64("test", "1234")
	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req = setHeaderContent(req, encodedCredentials)

	// Act
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

//
func (suite *TestSuite) Test_GetFeedMessages_Given_ValidUser_Returns_Ok() {
	// Arrange
	suite.registerUser("Testing", "tester@eh.com", "Testy")
	encodedCredentials := encodeCredentialsToB64("Testing", "Testy")
	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req = setHeaderContent(req, encodedCredentials)

	// Act
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// messages from self and followers
func (suite *TestSuite) Test_GetFeedMessages_Given_ValidUser_FollowingOtherUsers_returns_Messages() {
	// Arrange
	suite.registerUser("Testing", "tester@eh.com", "Testy")
	suite.registerUser("Jennifer", "jenni@eh.com", "Jesty")
	suite.registerUser("Geralt", "WhiteWolf@eh.com", "Gesty")

	testyCredentials := encodeCredentialsToB64("Testing", "Testy")
	gestyCredentials := encodeCredentialsToB64("Geralt", "Gesty")
	jestyCredentials := encodeCredentialsToB64("Jennifer", "Jesty")

	geraltMsg, _ := json.Marshal(gin.H{"content": "Gesty msg"})
	reqG := httptest.NewRequest(http.MethodPost, "/msgs/Geralt", bytes.NewReader(geraltMsg))
	reqG = setHeaderContent(reqG, gestyCredentials)
	suite.sendRequest(reqG)

	jenniferMsg, _ := json.Marshal(gin.H{"content": "Jesty msg"})
	reqJ := httptest.NewRequest(http.MethodPost, "/msgs/Jennifer", bytes.NewReader(jenniferMsg))
	reqJ = setHeaderContent(reqJ, jestyCredentials)
	suite.sendRequest(reqJ)

	time.Sleep(1 * time.Second)

	testyMsg, _ := json.Marshal(gin.H{"content": "Testy msg"})
	reqT := httptest.NewRequest(http.MethodPost, "/msgs/Testing", bytes.NewReader(testyMsg))
	reqT = setHeaderContent(reqT, testyCredentials)
	suite.sendRequest(reqT)

	fllwJennifer, _ := json.Marshal(gin.H{"follow": "Jennifer"})
	reqT = httptest.NewRequest(http.MethodPost, "/fllws/Testing", bytes.NewReader(fllwJennifer))
	reqT = setHeaderContent(reqT, testyCredentials)
	suite.sendRequest(reqT)

	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req = setHeaderContent(req, testyCredentials)

	msg0 := models.Message{Text: "Testy msg",}
	msg1 := models.Message{Text: "Jesty msg",}

	// Act
	w := suite.sendRequest(req)
	var resBody []models.Message
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), 2, len(resBody))
	assert.Equal(suite.T(), msg0.Text, resBody[0].Text)
	assert.Equal(suite.T(), msg1.Text, resBody[1].Text)
}

func (suite *TestSuite) registerUser(username, email, password string) {
	body, _ := json.Marshal(gin.H{"username": username, "email": email, "pwd": password})
	suite.sendRequest(httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body)))
}
