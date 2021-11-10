package models

import (
	"time"
)

type CSRFData struct {
	CSRFToken  string
	UserId     string
	Expiration time.Duration
}