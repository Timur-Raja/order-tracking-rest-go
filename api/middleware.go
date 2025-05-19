package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/user/usersql"
)

// Middleware is a collection of middleware functions for the Gin router

// ErrorLogger is a middleware that logs private errors
func ErrorLogger() gin.HandlerFunc {
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

// SessionAuth fetches useful information from the session which is accessing the API
func SessionAuth(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "authentication required"})
			return
		}

		// fetch the session
		query := usersql.NewSelectSessionByTokenQuery(db)
		query.Where.Token = token
		if err := query.Run(c); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				app.AbortWithErrorResponse(c, app.ErrAuthenticationRequired, err)
				return
			}
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
			return
		}
		//todo check expiry

		// save in context
		c.Set("user_id", query.Session.UserID)

		c.Next()
	}
}
