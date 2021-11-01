package csrf

import "errors"

var (
	ErrBadToken   = errors.New("bad token")
	ErrTokenNoExp = errors.New("token with no expiraton")
	ErrTokenExp   = errors.New("token has been expired")
)
