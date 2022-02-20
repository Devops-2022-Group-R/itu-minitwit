package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"

	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"

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

func (suite *MessageTestSuite) Test_GetFeedMessages_Given_NonExistent_User_Returns_NotFound() {
	// Act
	req, _ := http.NewRequest(http.MethodGet, "/feed", nil)
	w := suite.sendAuthedRequest(req, "test", "1234")

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

//
func (suite *MessageTestSuite) Test_GetFeedMessages_Given_ValidUser_Returns_Ok() {
	// Arrange
	suite.registerUser("Testing", "tester@eh.com", "Testy")

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/feed", nil)
	w := suite.sendAuthedRequest(req, "Testing", "Testy")

	// Assert
	suite.Equal(http.StatusOK, w.Code)
}

// messages from self and followers
func (suite *MessageTestSuite) Test_GetFeedMessages_Given_ValidUser_FollowingOtherUsers_returns_Messages() {
	// Arrange
	suite.registerUser("Testing", "tester@eh.com", "Testy")
	suite.registerUser("Jennifer", "jenni@eh.com", "Jesty")
	suite.registerUser("Geralt", "WhiteWolf@eh.com", "Gesty")

	geraltMsg, _ := json.Marshal(gin.H{"content": "Gesty msg"})
	reqG, _ := http.NewRequest(http.MethodPost, "/msgs/Geralt", bytes.NewReader(geraltMsg))
	suite.sendAuthedRequest(reqG, "Geralt", "Gesty")

	jenniferMsg, _ := json.Marshal(gin.H{"content": "Jesty msg"})
	reqJ, _ := http.NewRequest(http.MethodPost, "/msgs/Jennifer", bytes.NewReader(jenniferMsg))
	suite.sendAuthedRequest(reqJ, "Jennifer", "Jesty")

	time.Sleep(1 * time.Second)

	testyMsg, _ := json.Marshal(gin.H{"content": "Testy msg"})
	reqT, _ := http.NewRequest(http.MethodPost, "/msgs/Testing", bytes.NewReader(testyMsg))
	suite.sendAuthedRequest(reqT, "Testing", "Testy")

	fllwJennifer, _ := json.Marshal(gin.H{"follow": "Jennifer"})
	reqT, _ = http.NewRequest(http.MethodPost, "/fllws/Testing", bytes.NewReader(fllwJennifer))
	suite.sendAuthedRequest(reqT, "Testing", "Testy")

	msg0 := models.Message{Text: "Testy msg"}
	msg1 := models.Message{Text: "Jesty msg"}

	// Act
	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	w := suite.sendAuthedRequest(req, "Testing", "Testy")
	var resBody []models.Message
	suite.readBody(w, &resBody)

	// Assert
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(2, len(resBody))
	suite.Equal(msg0.Text, resBody[0].Text)
	suite.Equal(msg1.Text, resBody[1].Text)
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
