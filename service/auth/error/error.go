package error

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCookie       = errors.New("error with cookie")
	ErrEmptyData    = errors.New("required data is empty")
	ErrPostgres     = errors.New("internal DB server error")
)
