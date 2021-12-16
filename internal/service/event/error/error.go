package error

import "errors"

var (
	ErrEmptyData  = errors.New("required data is empty")
	ErrPostgres   = errors.New("internal DB server error")
	ErrAtoi       = errors.New("cant cast string to int")
	ErrNotAllowed = errors.New("user is not allowed to do this")
	ErrNoRows     = errors.New("no rows in a query result")
	ErrQuery      = errors.New("invalid query")
)
