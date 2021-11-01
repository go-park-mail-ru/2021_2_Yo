package auth

import "backend/models"

type UseCase interface {
	SignUp(user *models.User) (string, error)
	SignIn(mail, password string) (string, error)
	GetUser(userId string) (*models.User, error)
	UpdateUserInfo(userId string, name string, surname string, about string) error
	UpdateUserPassword(userId string, password string) error
}
