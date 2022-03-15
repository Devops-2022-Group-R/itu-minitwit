package internal

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// TODO: Move to a different file
type HttpError struct {
	StatusCode int
	Message    string
	Hidden     bool
	RelatedErr error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func NewHttpError(statusCode int, message string) HttpError {
	return HttpError{statusCode, message, false, nil}
}

func NewHiddenHttpError(statusCode int, message string) HttpError {
	return HttpError{statusCode, message, true, nil}
}

func NewHttpErrorWithRelatedError(statusCode int, message string, err error) HttpError {
	return HttpError{statusCode, message, false, err}
}

func NewHiddenHttpErrorWithRelatedError(statusCode int, message string, err error) HttpError {
	return HttpError{statusCode, message, true, err}
}

func NewInternalServerError(err error) HttpError {
	return NewHiddenHttpErrorWithRelatedError(500, "internal server error", err)
}

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}
