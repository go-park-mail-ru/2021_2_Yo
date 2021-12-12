package interfaces

import (
	authServiceModels "backend/internal/microservice/auth/models"
)

type SessionRepository interface {
	Create(data *authServiceModels.SessionData) error
	Check(sessionId string) (string, error)
	Delete(sessionId string) error
}
