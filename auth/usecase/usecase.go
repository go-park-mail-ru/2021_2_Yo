package usecase

import (
	"backend/auth"
	"backend/models"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

type UseCase struct {
	repository auth.Repository
	secretWord []byte
}

func NewUseCase(userRepo auth.Repository, secretWord []byte) *UseCase {
	return &UseCase{
		repository: userRepo,
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

func (a *UseCase) SignUp(name, surname, mail, password string) error {
	user := &models.User{
		Name:     name,
		Surname:  surname,
		Mail:     mail,
		Password: password,
	}
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail, password string) (string, error) {
	user, err := a.repository.GetUser(mail, password)
	if err == auth.ErrUserNotFound {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		ID: user.ID,
	})
	signedString, err := token.SignedString(a.secretWord)
	return signedString, err
}

func (a *UseCase) ParseToken(accessToken string) (*models.User, error) {
	userID, err := parseToken(accessToken, a.secretWord)
	if err != nil {
		return nil, err
	}
	return a.repository.GetUserById(userID)
}
