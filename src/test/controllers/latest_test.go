package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LatestTestSuite struct{ BaseTestSuite }

func TestLatestTestSuite(t *testing.T) {
	suite.Run(t, new(FollowTestSuite))
}

func (suite *LatestTestSuite) TestLatestController_SetLatest123_ReturnsLatest123() {
	// Arrange
	req1 := httptest.NewRequest(http.MethodGet, "/msgs?latest=123", nil)
	suite.sendRequest(req1)

	// Act
	req2 := httptest.NewRequest(http.MethodGet, "/latest", nil)
	w2 := suite.sendRequest(req2)

	var resBody map[string]int
	suite.readBody(w2, &resBody)

	// Assert
	suite.Equal(http.StatusOK, w2.Code)
	suite.Equal(123, resBody["latest"])
}

func (suite *LatestTestSuite) TestLatestController_SetLatestRepeated_ReturnsLastEachTime() {
	// Arrange
	doLatestCycle := func(value int) {
		// Arrange
		req1 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/latest?latest=%d", value), nil)
		suite.sendRequest(req1)

		// Act
		req2 := httptest.NewRequest(http.MethodGet, "/latest", nil)
		w2 := suite.sendRequest(req2)

		var resBody map[string]int
		suite.readBody(w2, &resBody)

		// Assert
		suite.Equal(http.StatusOK, w2.Code)
		suite.Equal(value, resBody["latest"])
	}

	// Act
	doLatestCycle(123)
	doLatestCycle(0)
	doLatestCycle(431)
	doLatestCycle(-1)
	doLatestCycle(1)
}

func (suite *LatestTestSuite) TestLatestController_GetLatestBeforeAnySet_Returns500() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/latest", nil)
	w := suite.sendRequest(req)

	// Act
	suite.Equal(http.StatusInternalServerError, w.Code)
}
