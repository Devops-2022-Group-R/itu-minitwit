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

func TestFlagToolTestSuite(t *testing.T) {
	suite.Run(t, new(FlagToolTestSuite))
}

func (suite *FlagToolTestSuite) Test_FlagMessageById_Given_hello_Returns_BadRequest() {

	// Act
	req := httptest.NewRequest(http.MethodPut, "/flag_tool/hello", nil)
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}


func (suite *FlagToolTestSuite) Test_GetAllMessages_Returns_AllMessages() {

	// 	Arrange
	const perPage = 30
	msg := "Gesty msg"
	user := "Geralt"
	suite.registerUser(user, "WhiteWolf@eh.com", "Gesty")
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
	req := httptest.NewRequest(http.MethodPut, "/flag_tool/msgs", nil)
	w := suite.sendRequest(req)

	var resBody []models.Message
	suite.readBody(w, &resBody)

	// 	Assert
	suite.Equal(http.StatusOK, w.Code)
	suite.True(len(resBody) > 31)
	suite.Equal(expectedMsg.Text, resBody[32].Text)
	suite.Equal(expectedMsg.Author.Username, resBody[32].Author.Username)
}
