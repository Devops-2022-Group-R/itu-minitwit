package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func NewBadRequestError(message string) HttpError {
	return NewHttpError(http.StatusBadRequest, message)
}

func NewBadRequestErrorFromError(err error) HttpError {
	return NewHttpErrorWithRelatedError(http.StatusBadRequest, err.Error(), err)
}

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

type returnedErr struct {
	Err        string `json:"error"`
	RelatedErr error  `json:"-"`
	Code       int    `json:"code"`
}

func ErrorHandleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		responseCode := 0

		if len(c.Errors) > 0 {
			errors := make([]returnedErr, 0)

			for _, err := range c.Errors {
				switch err.Err.(type) {
				case HttpError:
					httpErr := err.Err.(HttpError)
					if httpErr.RelatedErr != nil {
						Logger.Printf("[HTTP ERROR][CODE: %3d][%s] %s\n", httpErr.StatusCode, httpErr.Message, httpErr.RelatedErr.Error())
					} else {
						Logger.Printf("[HTTP ERROR][CODE: %3d] %s\n", httpErr.StatusCode, httpErr.Message)
					}

					if !httpErr.Hidden {
						if httpErr.StatusCode > responseCode {
							responseCode = httpErr.StatusCode
						}

						errors = append(errors, returnedErr{httpErr.Message, httpErr.RelatedErr, httpErr.StatusCode})
					} else {
						errors = append(errors, returnedErr{"Internal server error", nil, httpErr.StatusCode})
					}
				default:
					Logger.Printf("[INTERNAL ERROR][CODE: %3d] %s\n", 500, err.Err)
					responseCode = 500
					errors = append(errors, returnedErr{"Internal server error", nil, 500})
				}

			}

			c.JSON(responseCode, gin.H{
				"errors": errors,
			})
		}
	}
}
