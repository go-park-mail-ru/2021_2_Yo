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
}

func NewUseCaseAuth(userRepo auth.RepositoryUser) *UseCaseAuth {
	return &UseCaseAuth{
		userRepo:   userRepo,
		secretWord: []byte("welcometosanandreasimcjfromgrovestreet"),
	}
}

func (a *UseCaseAuth) SignUp(name, surname, mail, password string) error {
	user := &models.User{
		Name:     name,
		Surname:  surname,
		Mail:     mail,
		Password: password,
	}
	return a.userRepo.CreateUser(user)
}

func (a *UseCaseAuth) SignIn(mail, password string) (*models.User, string, error) {
	user, err := a.userRepo.GetUser(mail, password)
	if err == auth.ErrUserNotFound {
		return nil, "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		Name:    user.Name,
		Surname: user.Surname,
		Mail:    user.Mail,
	})
	signedString, err := token.SignedString(a.secretWord)
	return user, signedString, err
}

func (a *UseCaseAuth) List() []models.User {
	users := a.userRepo.List()
	var usersUsernames []models.User
	for _, user := range users {
		usersUsernames = append(usersUsernames, *user)
	}
	return usersUsernames
}

func (a *UseCaseAuth) ParseToken(cookie string) (string, error) {
	username, err := parser.ParseToken(cookie, a.secretWord)
	if err != nil {
		return "", err
	}
	return username, nil
}
