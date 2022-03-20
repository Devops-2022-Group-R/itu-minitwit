package internal

import "net/http"

var (
	ErrUserNotFound              = NewHttpError(http.StatusNotFound, "user not found")
	ErrUrlUsernameNotMatchHeader = NewHttpError(http.StatusUnprocessableEntity, "the URL username did not match the Authorization header username")
)
