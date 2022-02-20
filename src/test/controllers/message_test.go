package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func (suite *TestSuite) registerUser(username, email, password string) {
	body, _ := json.Marshal(gin.H{"username": username, "email": email, "pwd": password})
	suite.sendRequest(httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body)))
}
