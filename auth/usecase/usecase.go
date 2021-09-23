package usecase

import (
	"backend/auth"
	"backend/models"
)

type UseCaseAuth struct {
	userRepo auth.RepositoryUser
	//Добавить какую-нибудь информацию для токенов!
}

func NewUseCaseAuth(userRepo auth.RepositoryUser) *UseCaseAuth {
	return &UseCaseAuth{
		userRepo: userRepo,
	}
}

func (a *UseCaseAuth) SignUp(username, password string) error {
	user := &models.User{
		Username: username,
		Password: password,
	}
	return a.userRepo.CreateUser(user)
}

func (a *UseCaseAuth) SignIn(username, password string) (string, error) {
	user, err := a.userRepo.GetUser(username, password)
	if err != nil {
		return "", err
	}
	//Здесь нужна работа с токенами!
	return user.Username, nil
}
