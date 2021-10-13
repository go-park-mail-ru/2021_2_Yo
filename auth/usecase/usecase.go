package usecase

import (
	"backend/auth"
	"backend/models"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	log "github.com/sirupsen/logrus"
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
	password_hash := a.CreatePasswordHash(password)
	user := &models.User{
		Name:     name,
		Surname:  surname,
		Mail:     mail,
		Password: password_hash,
	}
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail, password string) (string, error) {
	password_hash := a.CreatePasswordHash(password)
	log.Info("UseCase:SignIn password:", password, " password_hash:", password_hash)
	user, err := a.repository.GetUser(mail, password_hash)
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

func (a *UseCase) CreatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
