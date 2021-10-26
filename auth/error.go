package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrInvalidUserData    = errors.New("invalid user data")
	ErrUserExists         = errors.New("user with the same mail already exists")
	ErrExpired            = errors.New("need to refresh token")
)
