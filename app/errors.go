package app

import "errors"

var (
	ErrNoEnvFileFound   = errors.New("no ENV file found")
	ErrNoEnvKeyProvided = errors.New("no ENV key provided")
	ErrNoEnvValueFound  = errors.New("no ENV value found")
)
