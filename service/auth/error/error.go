package error

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCookie       = errors.New("error with cookie")
	ErrPostgres     = errors.New("internal DB server error")
	ErrUserExists   = errors.New("user already exists")
)
