package auth

import (
	"backend/models"
)

type Repository interface {
	GetUserById(userId string) (*models.User, error)
	UpdateUserInfo(userId, name, surname, about string) error
	UpdateUserPassword(userId, password string) error
}
