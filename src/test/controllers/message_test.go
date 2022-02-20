package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MessageTestSuite) Test_GetUserMessages_Returns_OK() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")

	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MessageTestSuite) Test_GetUserMessages_Given_NonExistent_User_Returns_NotFound() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/msgs/Darrow", nil)
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_Returns_NoContent() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")

	// Act
	body, _ := json.Marshal(gin.H{"content": "Omnis vir lupus."})
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendAuthedRequest(req, "Darrow", "Reaper")

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_Given_NonExistent_User_Returns_NotFound() {
	// Act
	body, _ := json.Marshal(gin.H{"content": "Omnis vir lupus."})
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendAuthedRequest(req, "Darrow", "Reaper")

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_Given_Empty_Message_Returns_BadRequest() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")

	// Act
	body, _ := json.Marshal(gin.H{"content": ""})
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendAuthedRequest(req, "Darrow", "Reaper")

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_As_Simulator_Returns_No_Content() {
	// Arrange
	suite.registerUser("Darrow", "darrow@andromedus.com", "Reaper")
	suite.registerSimulator()

	// Act
	body, _ := json.Marshal(gin.H{"content": "Some message"})
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *MessageTestSuite) Test_PostUserMessage_As_Simulator_With_Unknown_User_Returns_Not_Found() {
	// Arrange
	suite.registerSimulator()

	// Act
	body, _ := json.Marshal(gin.H{"content": "Some message"})
	req := httptest.NewRequest(http.MethodPost, "/msgs/Darrow", bytes.NewReader(body))
	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}
