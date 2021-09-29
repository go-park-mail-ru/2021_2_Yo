package usecase

import (
	"backend/auth"
	"backend/models"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

type UseCaseAuth struct {
	userRepo   auth.RepositoryUser
	secretWord []byte
}

func NewUseCaseAuth(userRepo auth.RepositoryUser, secretWord []byte) *UseCaseAuth {
	return &UseCaseAuth{
		userRepo:   userRepo,
		secretWord: secretWord,
	}
}

type claims struct {
	jwt.StandardClaims
	ID string `json:"user_id"`
}

func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.ID, nil
	}
	return "", auth.ErrInvalidAccessToken
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

func (a *UseCaseAuth) SignIn(mail, password string) (string, error) {
	user, err := a.userRepo.GetUser(mail, password)
	if err == auth.ErrUserNotFound {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		ID: user.ID,
	})
	signedString, err := token.SignedString(a.secretWord)
	return signedString, err
}

func (a *UseCaseAuth) ParseToken(cookie string) (string, error) {
	userID, err := parseToken(cookie, a.secretWord)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (a *UseCaseAuth) GetUserById(userID string) (*models.User, error) {
	return a.userRepo.GetUserById(userID)
}
