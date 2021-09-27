package auth

import "backend/models"

type UseCase interface {
	SignUp(name, surname, mail, password string) error
	SignIn(mail, password string) (string, error)
	ParseToken(cookie string) (string, error)
	List() []models.User
}
