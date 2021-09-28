package auth

import "backend/models"

type UseCaseAuth interface {
	SignUp(name, surname, mail, password string) error
	SignIn(mail, password string) (string, error)
	ParseToken(cookie string) (string, error)
	GetUserById(userID string) (*models.User, error)
}
