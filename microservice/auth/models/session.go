package models

import (
	"time"
)

type SessionData struct {
	SessionId  string
	UserId     string
	Expiration time.Duration
}