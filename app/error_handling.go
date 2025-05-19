package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// custom error type to streamline error handling
type apiError struct {
	Status int
	Err    error
}

func NewAPIError(status int, desc string) *apiError {
	return &apiError{Status: status, Err: errors.New(desc)}
}

// send error as private for internal logging and abort call with status and error message
func AbortWithErrorResponse(c *gin.Context, apiError *apiError, privateError ...error) {
	// if we have privateErrors log them, otherwise log the custom error
	if len(privateError) > 0 {
		for _, err := range privateError {
			c.Error(err).SetType(gin.ErrorTypePrivate)
		}
	} else {
		c.Error(apiError.Err).SetType(gin.ErrorTypePrivate)
	}

	c.AbortWithStatusJSON(apiError.Status, gin.H{"error": apiError.Err.Error()})
}

// generic errors
var (
	ErrFailedToLoadParams     = NewAPIError(http.StatusInternalServerError, "failed to load params")
	ErrMissingBodyParams      = NewAPIError(http.StatusBadRequest, "missing request body params")
	ErrServerError            = NewAPIError(http.StatusInternalServerError, "internal server error")
	ErrAuthenticationRequired = NewAPIError(http.StatusUnauthorized, "Authentication required")
	ErrResourceNotFound       = NewAPIError(http.StatusNotFound, "resource not found")
)
