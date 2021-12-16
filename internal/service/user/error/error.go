package error

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyData    = errors.New("required data is empty")
	ErrPostgres     = errors.New("internal DB server error")
	ErrAtoi         = errors.New("cant cast string to int")
)
