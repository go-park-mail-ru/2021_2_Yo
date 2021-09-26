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

func (a *UseCaseAuth) SignUp(name, surname, mail,password string) error {
	user := &models.User{
		Name: name,
		Surname: surname,
		Mail: mail,
		Password: password,
	}
	return a.userRepo.CreateUser(user)
}

func (a *UseCaseAuth) SignIn(mail,password string) (string, error) {
	user, err := a.userRepo.GetUser(mail,password)
	if err == auth.ErrUserNotFound {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		Name: user.Name,
		Surname: user.Surname,
		Mail: user.Mail,
		
	})
	//return user.Username, nil
	return token.SignedString(a.secretWord)
}

func (a *UseCaseAuth) List() []string {
	users := a.userRepo.List()
	var usersUsernames []string
	for _, user := range users {
		usersUsernames = append(usersUsernames, user.Mail)
	}
	return usersUsernames
}

func (a *UseCaseAuth) Parse(cookie string) (string, error) {
	username, err := parser.ParseToken(cookie, a.secretWord)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (a *UseCaseAuth) ParseKsenia(cookie string) (string, error) {
	username, err := parser.ParseTokenForKsenia(cookie, a.secretWord)
	if err != nil {
		return "", err
	}
	return username, nil
}