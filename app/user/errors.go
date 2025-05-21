package user

import (
	"net/http"

	"github.com/timur-raja/order-tracking-rest-go/app"
)

var (
	ErrUserNotFound       = app.NewAPIError(http.StatusNotFound, "user not found")
	ErrInvalidCredentials = app.NewAPIError(http.StatusUnauthorized, "invalid credentials")
	ErrUserAlreadyExists  = app.NewAPIError(http.StatusConflict, "user already exists")
)
