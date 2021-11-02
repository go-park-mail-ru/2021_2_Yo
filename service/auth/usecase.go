package auth

import "backend/models"

type UseCase interface {
	SignUp(user *models.User) (string, error)
	SignIn(mail, password string) (string, error)
}
