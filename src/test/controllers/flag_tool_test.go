package controllers_test

import (
	_ "bytes"
	_ "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FlagToolTestSuite struct{ BaseTestSuite }

func TestFlagToolTestSuite(t *testing.T) {
	suite.Run(t, new(FlagToolTestSuite))
}

func (suite *FlagToolTestSuite) Test_FlagMessageById_Given_hello_Returns_BadRequest() {
	// Act
	req := httptest.NewRequest(http.MethodGet, "/flag_tool/hello", nil)
	w := suite.sendRequest(req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}
