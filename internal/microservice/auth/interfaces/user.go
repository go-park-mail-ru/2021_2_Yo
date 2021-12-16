package interfaces

import (
	"backend/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (string, error)
	GetUser(mail, password string) (*models.User, error)
}
