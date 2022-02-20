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

type FollowTestSuite struct{ BaseTestSuite }

func TestFollowTestSuite(t *testing.T) {
	suite.Run(t, new(FollowTestSuite))
}

const followUrl = "/fllws/geralt"

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
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w1 := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})
	w2 := sendAuthRequest(suite, followUrl, gin.H{"follow": "triss"})

	// Assert
	suite.Equal(http.StatusNoContent, w1.Code)
	suite.Equal(http.StatusNoContent, w2.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidFollow_Returns204() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenNonExistingUser_Returns404() {
	// Act
	w := sendAuthRequest(suite, "/fllws/i-dont-exist", gin.H{"follow": "yennefer"})

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenFollowANonExistingUser_Returns404() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "vesemir"})

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenWrongUsernameInURL_Returns422() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, "/fllws/eredin", gin.H{"follow": "yennefer"})

	// Assert
	suite.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_WithoutLoggedInUser_Returns403() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	body, _ := json.Marshal(gin.H{"follow": "yennefer"})
	req := httptest.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *FollowTestSuite) TestFollowGetController_GivenUserWithNoFollows_ReturnsEmpty() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// act
	req, _ := http.NewRequest(http.MethodGet, followUrl, nil)
	w := suite.sendRequest(req)

	var resBody []string
	suite.readBody(w, &resBody)

	// Asssert
	suite.Empty(resBody)
}

func (suite *FollowTestSuite) TestFollowGetController_GivenUserWithFollowed_ReturnsAllFollowed() {
	// Arrange
	setupTestFollowRelationships(suite)

	// Act
	req, _ := http.NewRequest(http.MethodGet, followUrl, nil)
	w := suite.sendRequest(req)

	var resBody map[string][]string
	suite.readBody(w, &resBody)

	// Assert
	suite.Equal(http.StatusOK, w.Code)
	suite.ElementsMatch([...]string{"yennefer", "triss"}, resBody["follows"])
}

func (suite *FollowTestSuite) TestFollowPostController_GivenUnfollowANonExistingUser_Returns404() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "vesemir"})

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidUnfollow_Returns204() {
	// Arrange
	setupTestFollowRelationships(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss"})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenFollowAlreadyFollowed_Returns204() {
	// Arrange
	setupTestFollowRelationships(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"follow": "yennefer"})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenUnfollowAlreadyUnfollowed_Returns204() {
	// Arrange
	setupTestFollowRelationships(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "eredin"})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenValidFollowAndUnfollow_Returns204() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{"unfollow": "triss", "follow": "eredin"})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_GivenEmptyBody_Returns204() {
	// Arrange
	utilCreateUsersInDatabase(suite)

	// Act
	w := sendAuthRequest(suite, followUrl, gin.H{})

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenValidFollow_Returns204() {
	// Arrange
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	// Act
	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenValidUnfollow_Returns204() {
	// Arrange
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	// Act
	body, _ := json.Marshal(gin.H{"unfollow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, followUrl, bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *FollowTestSuite) TestFollowPostController_AsAdminGivenNonExistingUser_Returns404() {
	// Arrange
	setupTestFollowRelationships(suite)
	suite.registerSimulator()

	// Act
	body, _ := json.Marshal(gin.H{"follow": "triss"})
	req, _ := http.NewRequest(http.MethodPost, "/fllws/eskel", bytes.NewReader(body))

	w := suite.sendSimulatorRequest(req)

	// Assert
	suite.Equal(http.StatusNotFound, w.Code)
}
