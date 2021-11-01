package event

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEmptyData     = errors.New("required data is empty")
	ErrPostgres      = errors.New("internal DB server error")
	ErrAtoi          = errors.New("cant cast string to int")
	ErrNotAllowed    = errors.New("user is not allowed to do this")
)
