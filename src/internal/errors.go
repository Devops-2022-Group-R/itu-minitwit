package internal

import "net/http"

var (
	ErrUserNotFound = NewHttpError(http.StatusNotFound, "user not found")
)
