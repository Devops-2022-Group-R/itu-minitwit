package controllers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterController_GivenNoBody_Returns400(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/register", nil)
	w := sendRequest(req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
