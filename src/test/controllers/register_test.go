package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestRegisterController_GivenNoBody_Returns400() {
	req, _ := http.NewRequest(http.MethodPost, "/register", nil)
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *TestSuite) TestRegisterController_GivenValidRequest_Returns204() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yennefer@aretuza.wr",
		"password": "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *TestSuite) TestRegisterController_GivenInvalidEmail_Returns422() {
	body, _ := json.Marshal(gin.H{
		"username": "Yennefer of Vengerberg",
		"email":    "yenneferaretuza.wr",
		"password": "chaosmaster",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := suite.sendRequest(req)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
}
