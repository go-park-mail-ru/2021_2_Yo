package auth

import (
	"backend/models"
)

type Repository interface {
	CreateUser(user *models.User) (string, error)
	GetUser(mail, password string) (*models.User, error)
	GetUserById(userId string) (*models.User, error)
}
