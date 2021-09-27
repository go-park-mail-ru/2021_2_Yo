package auth

import "backend/models"

type UseCase interface {
	SignUp(name, surname, mail, password string) error
	SignIn(mail, password string) (string, error)
	ParseToken(cookie string) (string, error)
	//КОСТЫЛИ!!!
	GetUser(mail, password string) (*models.User, error)
	GetUserById(userID string) (*models.User, error)
}
