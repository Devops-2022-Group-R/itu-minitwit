package controllers_test

import (
	"encoding/base64"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestLoginController_GivenNoHeader_Returns401() {
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_Bruce_Lee_10000kicks_returns404() {
	req, _ := http.NewRequest(http.MethodPost, "/login", nil)
	encodedCredentials := encodeCredentialsToB64("Bruce Lee", "10000kicks")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_ValidUsersCredentials_returns204() {
	addUserToTestDb(suite)
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	encodedCredentials := encodeCredentialsToB64("Yennefer of V", "chaosmaster")
	req = setHeaderContent(req, encodedCredentials)

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestLoginController_Given_ValidUserAndInvalidPassword_returns401() {
	addUserToTestDb(suite)
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	encodedCredentials := encodeCredentialsToB64("Yennefer of V", "chaos")
	req = setHeaderContent(req, encodedCredentials)

	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}
func encodeCredentialsToB64(username string, password string) string {
	data := username + ":" + password
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}
