package usecase

import (
	"backend/auth"
	"backend/models"
	//"time"
	"github.com/dgrijalva/jwt-go/v4"
	//"crypto/sha1"
	"backend/parser"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		Username: user.Username,
	})
	//return user.Username, nil
	return token.SignedString(a.secretWord)
}

func (a *UseCaseAuth) List() []string {
	users := a.userRepo.List()
	var usersUsernames []string
	for _, user := range users {
		usersUsernames = append(usersUsernames, user.Username)
	}
	return usersUsernames
}

func (a *UseCaseAuth) Parse(cookie string) (string, error) {
	username, err := parser.ParseToken(cookie, a.secretWord)
	if err != nil {
		return username, nil
	}
	return "", err
}
