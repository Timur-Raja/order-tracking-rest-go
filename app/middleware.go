package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Middleware is a collection of middleware functions for the Gin router

func ErrorLogger() gin.HandlerFunc {
	// ErrorLogger is a middleware that logs private errors
	return func(c *gin.Context) {
		c.Next()

		// detect eventual private errors and log them
		if len(c.Errors.ByType(gin.ErrorTypePrivate)) > 0 {
			err := c.Errors.ByType(gin.ErrorTypePrivate)[0].Err
			logrus.WithError(err).
				WithField("path", c.FullPath()).
				Error("request failed")
		}
	}
}
