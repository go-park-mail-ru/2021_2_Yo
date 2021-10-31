package auth

import "backend/models"

type UseCase interface {
	SignUp(user *models.User) error
	SignIn(mail, password string) (string, error)
	ParseToken(accessToken string) (*models.User, error)
	Logout(accessToken string) (string, error)
	GetCSRFToken(cookie string, expirationTime int64) (string, error)
}
