package user

import (
	"backend/models"
)

type Repository interface {
	GetUserById(userId string) (*models.User, error)
	UpdateUserInfo(user *models.User) error
	UpdateUserPassword(userId, password string) error
}
