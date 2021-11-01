package error

import "errors"

var (
	ErrTokenNoExp 		= errors.New("token with no expiraton")
	ErrTokenExp   		= errors.New("token has been expired")
	ErrEmptyToken 	    = errors.New("token is empty")
	ErrRedis          	= errors.New("internal redis error")
)
