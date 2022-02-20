package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const followUrl = "/fllws/geralt"

// Unsafe to use unless you are strictly using the username
func utilCreateUsersInDatabase(suite *TestSuite) {
	suite.registerUser("geralt", "geralt@witcher.pl", "123")
	suite.registerUser("yennefer", "yennefer@witcher.pl", "123")
	suite.registerUser("triss", "triss@witcher.pl", "123")
	suite.registerUser("eredin", "eredin@wildhunt.pl", "123")
}

func setupTestFollowRelationships(suite *TestSuite) {
	utilCreateUsersInDatabase(suite)
	w1 := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	w2 := sendAuthRequest(suite, followUrl, gin.H{"follow": "triss"})
	assert.Equal(suite.T(), http.StatusNoContent, w1.Code)
	assert.Equal(suite.T(), http.StatusNoContent, w2.Code)
}

func sendAuthRequest(suite *TestSuite, url string, body gin.H) *httptest.ResponseRecorder {
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Basic "+encodeCredentialsToB64("geralt", "123"))
	return suite.sendRequest(req)
}

func (suite *TestSuite) TestFollowPostController_GivenValidFollow_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenNonExistingUser_Returns404() {
	w := sendAuthRequest(suite, "/fllws/i-dont-exist", gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenFollowANonExistingUser_Returns404() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "vesemir"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenWrongUsernameInURL_Returns422() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, "/fllws/eredin", gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
}

func (suite *TestSuite) TestFollowPostController_WithoutLoggedInUser_Returns403() {
	utilCreateUsersInDatabase(suite)

	body, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req := httptest.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *TestSuite) TestFollowGetController_GivenUserWithNoFollows_ReturnsEmpty() {
	utilCreateUsersInDatabase(suite)

	req, _ := http.NewRequest(http.MethodGet, followUrl, nil)
	w := suite.sendRequest(req)

	var resBody []string
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Empty(suite.T(), resBody)
}

func (suite *TestSuite) TestFollowGetController_GivenUserWithFollowed_ReturnsAllFollowed() {
	setupTestFollowRelationships(suite)
	assert := assert.New(suite.T())

	req, _ := http.NewRequest(http.MethodGet, followUrl, nil)
	w := suite.sendRequest(req)

	var resBody map[string][]string
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Equal(http.StatusOK, w.Code)
	assert.ElementsMatch([...]string{"yennefer", "triss"}, resBody["follows"])
}

func (suite *TestSuite) TestFollowPostController_GivenUnfollowANonExistingUser_Returns404() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "vesemir"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenValidUnfollow_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenFollowAlreadyFollowed_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenUnfollowAlreadyUnfollowed_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "eredin"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenValidFollowAndUnfollow_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss", "follow": "eredin"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_GivenEmptyBody_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_AsAdminGivenValidFollow_Returns204() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_AsAdminGivenValidUnfollow_Returns204() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"unfollow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestFollowPostController_AsAdminGivenNonExistingUser_Returns404() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, "/fllws/eskel", bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}
