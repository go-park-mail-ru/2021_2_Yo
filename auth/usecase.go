package auth

import "backend/models"

type UseCase interface {
	SignUp(name, surname, mail, password string) error
	SignIn(mail, password string) (*models.User, string, error)
	ParseToken(cookie string) (string, error)
	List() []models.User
}
