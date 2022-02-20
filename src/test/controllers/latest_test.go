package controllers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestLatestController_SetLatest123_ReturnsLatest123() {
	req1 := httptest.NewRequest(http.MethodGet, "/msgs?latest=123", nil)
	// We do not care about the result of this request
	suite.sendRequest(req1)

	req2 := httptest.NewRequest(http.MethodGet, "/latest", nil)
	w2 := suite.sendRequest(req2)

	var resBody map[string]int
	resBodyBytes, _ := ioutil.ReadAll(w2.Result().Body)
	json.Unmarshal(resBodyBytes, &resBody)

	assert.Equal(suite.T(), http.StatusOK, w2.Code)
	assert.Equal(suite.T(), 123, resBody["latest"])
}

func (suite *TestSuite) TestLatestController_SetLatestRepeated_ReturnsLastEachTime() {
	doLatestCycle := func(value int) {
		req1 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/latest?latest=%d", value), nil)
		suite.sendRequest(req1)
		req2 := httptest.NewRequest(http.MethodGet, "/latest", nil)
		w2 := suite.sendRequest(req2)
		var resBody map[string]int
		resBodyBytes, _ := ioutil.ReadAll(w2.Result().Body)
		json.Unmarshal(resBodyBytes, &resBody)
		assert.Equal(suite.T(), http.StatusOK, w2.Code)
		assert.Equal(suite.T(), value, resBody["latest"])
	}

	doLatestCycle(123)
	doLatestCycle(0)
	doLatestCycle(431)
	doLatestCycle(-1)
	doLatestCycle(1)
}

func (suite *TestSuite) TestLatestController_GetLatestBeforeAnySet_Returns500() {
	req := httptest.NewRequest(http.MethodGet, "/latest", nil)
	w := suite.sendRequest(req)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
