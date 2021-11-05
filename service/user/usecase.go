package user

import "backend/models"

type UseCase interface {
	GetUserById(userId string) (*models.User, error)
	UpdateUserInfo(userId string, name string, surname string, about string) error
	UpdateUserPassword(userId string, password string) error
}
