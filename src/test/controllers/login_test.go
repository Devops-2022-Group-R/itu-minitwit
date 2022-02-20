package controllers_test

import (
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestLoginController_GivenNoHeader_Returns401() {
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_UnknownUser_returns404() {
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Bruce Lee", "10000kicks"))

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_ValidUsersCredentials_returns204() {
	suite.registerUser("Yennefer of V", "yennefer@aretuza.wr", "chaosmaster")

	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("Yennefer of V", "chaosmaster"))

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_ValidUserAndInvalidPassword_returns401() {
	suite.registerUser("Yennefer of V", "yennefer@aretuza.wr", "chaosmaster")

	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", encodeCredentialsToB64("Yennefer of V", "chaos"))

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}
