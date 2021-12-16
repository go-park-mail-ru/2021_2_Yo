package auth

import (
	"backend/internal/models"
)

type UseCase interface {
	SignUp(u *models.User) (string, error)
	SignIn(u *models.User) (string, error)
	CreateSession(userId string) (string, error)
	CheckSession(SessionId string) (string, error)
	DeleteSession(SessionId string) error
	CreateToken(userId string) (string, error)
	CheckToken(csrfToken string) (string, error)
}
