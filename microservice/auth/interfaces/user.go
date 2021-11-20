package interfaces

import (
	"backend/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (string, error)
	GetUser(mail, password string) (*models.User, error)
}