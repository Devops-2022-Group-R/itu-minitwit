package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct{ BaseTestSuite }

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

func (suite *MessageTestSuite) Test_GetMessages_Returns_OK() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *MessageTestSuite) Test_GetUserMessages_Returns_OK() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")

	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *MessageTestSuite) Test_GetUserMessages_Given_NonExistent_User_Returns_NotFound() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_Returns_NoContent() {
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

func (suite *MessageTestSuite) Test_PostUserMessage_Given_NonExistent_User_Returns_NotFound() {
	// Arrange
	body, _ := json.Marshal(gin.H{"content": "Omnis vir lupus."})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Darrow", "Reaper"))
	w := suite.sendRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_Given_Empty_Message_Returns_BadRequest() {
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

func (suite *MessageTestSuite) Test_PostUserMessage_As_Simulator_Returns_No_Content() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"content": "Some message"})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendSimulatorRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_As_Simulator_With_Unknown_User_Returns_Not_Found() {
	// Arrange
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"content": "Some message"})

	// Act
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendSimulatorRequest(req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}
