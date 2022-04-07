package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FlagToolTestSuite struct{ BaseTestSuite }

const (
	user                = "Geralt"
	flagToolMsgEndpoint = "/flag_tool/msgs"
)

func TestFlagToolTestSuite(t *testing.T) {
	suite.Run(t, new(FlagToolTestSuite))
}

func utilCreateUserInDatabase(suite *FlagToolTestSuite) {
	suite.registerUser(user, "WhiteWolf@eh.com", "Gesty")
}

func (suite *FlagToolTestSuite) Test_FlagMessageById_Given_hello_Returns_BadRequest() {

	// Act
	req := httptest.NewRequest(http.MethodPut, "/flag_tool/hello", nil)
	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *FlagToolTestSuite) Test_FlagMessageById_Given_ExistingMsgId_Returns_Message_Flagged() {

	// Arrange
	msg := "Gesty msg"
	utilCreateUserInDatabase(suite)

	geraltMsg, _ := json.Marshal(gin.H{"content": msg})
	reqG, _ := http.NewRequest(http.MethodPost, "/msgs/Geralt", bytes.NewReader(geraltMsg))
	suite.sendAuthedRequest(reqG, user, "Gesty")

	expectedMsg := models.Message{
		Author: models.User{
			Username: user},
		Text:    msg,
		Flagged: true,
	}

	// Act
	reqPreFlag := httptest.NewRequest(http.MethodGet, flagToolMsgEndpoint, nil)
	wPreFlag := suite.sendSimulatorRequest(reqPreFlag)

	var resBodyPreFlag []models.Message
	suite.readBody(wPreFlag, &resBodyPreFlag)

	req := httptest.NewRequest(http.MethodPut, "/flag_tool/1", nil)
	w := suite.sendSimulatorRequest(req)

	reqPostFlag := httptest.NewRequest(http.MethodGet, flagToolMsgEndpoint, nil)
	wPostFlag := suite.sendSimulatorRequest(reqPostFlag)

	var resBodyPostFlag []models.Message
	suite.readBody(wPostFlag, &resBodyPostFlag)

	// Assert
	// pre flag
	suite.Equal(http.StatusOK, wPreFlag.Code)
	suite.False(resBodyPreFlag[0].Flagged)
	suite.Equal(expectedMsg.Text, resBodyPreFlag[0].Text)
	suite.Equal(expectedMsg.Author.Username, resBodyPreFlag[0].Author.Username)

	//	post flag
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(http.StatusOK, wPostFlag.Code)
	suite.Equal(expectedMsg.Text, resBodyPostFlag[0].Text)
	suite.Equal(expectedMsg.Author.Username, resBodyPostFlag[0].Author.Username)
	suite.True(resBodyPostFlag[0].Flagged)
}

func (suite *FlagToolTestSuite) Test_GetAllMessages_Returns_AllMessages() {

	// 	Arrange
	const perPage = 30
	msg := "Gesty msg"
	utilCreateUserInDatabase(suite)

	// add more messages then perpage limit
	for i := 0; i < perPage+5; i++ {
		geraltMsg, _ := json.Marshal(gin.H{"content": msg})
		reqG, _ := http.NewRequest(http.MethodPost, "/msgs/Geralt", bytes.NewReader(geraltMsg))
		suite.sendAuthedRequest(reqG, user, "Gesty")
	}

	expectedMsg := models.Message{
		Author: models.User{
			Username: user},
		Text: msg,
	}

	// 	Act
	req := httptest.NewRequest(http.MethodGet, flagToolMsgEndpoint, nil)
	w := suite.sendSimulatorRequest(req)

	var resBody []models.Message
	suite.readBody(w, &resBody)

	// 	Assert
	suite.Equal(http.StatusOK, w.Code)
	suite.True(len(resBody) > 31)
	suite.Equal(expectedMsg.Text, resBody[32].Text)
	suite.Equal(expectedMsg.Author.Username, resBody[32].Author.Username)
}

func (suite *FlagToolTestSuite) Test_GetAllMessages_Without_Authorization_Returns_ForbiddenRequest() {

	// Arrange
	utilCreateUserInDatabase(suite)

	// Act
	req := httptest.NewRequest(http.MethodGet, flagToolMsgEndpoint, nil)
	w := suite.sendAuthedRequest(req, user, "Gesty")

	// Assert
	suite.Equal(http.StatusForbidden, w.Code)
}
