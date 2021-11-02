package error

import "errors"

var (
	ErrEmptySessionId = errors.New("session id is empty")
	ErrRedis          = errors.New("internal redis error")
	ErrCreateSession  = errors.New("session was not created")
	ErrCheckSession   = errors.New("session wasn't got")
	ErrDeleteSession  = errors.New("session was not deleted")
)
