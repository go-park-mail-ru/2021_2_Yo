package interfaces

import (
	"backend/pkg/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (string, error)
	GetUser(mail, password string) (*models.User, error)
}
