package user

import "errors"

var (
	ErrMissingParams = errors.New("body parameters are required")
)
