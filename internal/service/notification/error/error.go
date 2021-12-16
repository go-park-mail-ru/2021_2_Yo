package error

import "errors"

var (
	ErrPostgres = errors.New("internal DB server error")
	ErrNoRows   = errors.New("no rows in a query result")
)
