package interfaces

import (
	"backend/service/session/models"
)

type SessionRepository interface {
	Create(data *models.SessionData) error 
	Check(sessionId string) (string, error)
	Delete(sessionId string) error
}


