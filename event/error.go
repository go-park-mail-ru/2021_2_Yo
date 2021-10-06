package event

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event is already exists")
)
