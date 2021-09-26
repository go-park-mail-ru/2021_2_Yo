package usecase

import (
	"backend/auth"
	"backend/models"
	//"crypto/sha1"
)

type UseCaseAuth struct {
	userRepo   auth.RepositoryUser
	secretWord []byte
	//Добавить какую-нибудь информацию для токенов!
}

func NewUseCaseAuth(userRepo auth.RepositoryUser) *UseCaseAuth {
	return &UseCaseAuth{
		userRepo:   userRepo,
		secretWord: []byte("welcometosanandreasimcjfromgrovestreet"),
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
	if err == auth.ErrUserNotFound {
		return "", err
	}
	//Здесь нужна работа с токенами!
	return user.Username, nil
}

func (a *UseCaseAuth) List() []string{
	users := a.userRepo.List()
	var usersUsernames []string
	for _,user := range users {
		usersUsernames = append(usersUsernames, user.Username)
	}
	return usersUsernames
}