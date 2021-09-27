package auth

import (
	"backend/models"
)

type RepositoryUser interface {
	CreateUser(user *models.User) error
	GetUser(mail, password string) (*models.User, error)
	List() []*models.User
}
