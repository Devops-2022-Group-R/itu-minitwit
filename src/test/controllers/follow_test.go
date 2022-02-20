package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const followUrl = "/fllws/geralt"

type FollowTestSuite struct{ BaseTestSuite }

func TestFollowTestSuite(t *testing.T) {
	suite.Run(t, new(FollowTestSuite))
}

// Unsafe to use unless you are strictly using the username
func utilCreateUsersInDatabase(suite *FollowTestSuite) {
	suite.registerUser("geralt", "geralt@witcher.pl", "123")
	suite.registerUser("yennefer", "yennefer@witcher.pl", "123")
	suite.registerUser("triss", "triss@witcher.pl", "123")
	suite.registerUser("eredin", "eredin@wildhunt.pl", "123")
}

func sendAuthRequest(suite *FollowTestSuite, url string, body gin.H) *httptest.ResponseRecorder {
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))

	return suite.sendAuthedRequest(req, "geralt", "123")
}

func setupTestFollowRelationships(suite *FollowTestSuite) {
	utilCreateUsersInDatabase(suite)
	w1 := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	w2 := sendAuthRequest(suite, followUrl, gin.H{"follow": "triss"})
	assert.Equal(suite.T(), http.StatusNoContent, w1.Code)
	assert.Equal(suite.T(), http.StatusNoContent, w2.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidFollow_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenNonExistingUser_Returns404() {
	w := sendAuthRequest(suite, "/fllws/i-dont-exist", gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenFollowANonExistingUser_Returns404() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "vesemir"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenWrongUsernameInURL_Returns422() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, "/fllws/eredin", gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_WithoutLoggedInUser_Returns403() {
	utilCreateUsersInDatabase(suite)

	body, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req := httptest.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *FollowTestSuite) TestFollowGetController_GivenUserWithNoFollows_ReturnsEmpty() {
	utilCreateUsersInDatabase(suite)

	req, _ := http.NewRequest(http.MethodGet, followUrl, nil)
	w := suite.sendRequest(req)

	var resBody []string
	resBodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Empty(suite.T(), resBody)
}

func (suite *FollowTestSuite) TestFollowGetController_GivenUserWithFollowed_ReturnsAllFollowed() {
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

func (suite *FollowTestSuite) TestFollowPostController_GivenUnfollowANonExistingUser_Returns404() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "vesemir"})
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidUnfollow_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenFollowAlreadyFollowed_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenUnfollowAlreadyUnfollowed_Returns204() {
	setupTestFollowRelationships(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "eredin"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidFollowAndUnfollow_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss", "follow": "eredin"})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenEmptyBody_Returns204() {
	utilCreateUsersInDatabase(suite)
	w := sendAuthRequest(suite, followUrl, gin.H{})
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenValidFollow_Returns204() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenValidUnfollow_Returns204() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"unfollow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenNonExistingUser_Returns404() {
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, "/fllws/eskel", bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}
