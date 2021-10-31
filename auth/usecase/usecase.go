package usecase

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
	"backend/csrf"
)

const logMessage = "auth:usecase:usecase:"

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
		if claims.ExpiresAt.Before(time.Now()) {
			log.Debug("Expired")
			return "", auth.ErrExpired
		}
		return claims.ID, nil
	}
	return "", auth.ErrInvalidAccessToken
}

func (a *UseCase) SignUp(user *models.User) error {
	password_hash := a.CreatePasswordHash(user.Password)
	user.Password = password_hash
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail, password string) (string, error) {
	message := logMessage + "SignIn:"
	hashedPassword := a.CreatePasswordHash(password)
	log.Debug(message+"hashedPassword =", hashedPassword)
	user, err := a.repository.GetUser(mail, hashedPassword)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		jwt.StandardClaims{ExpiresAt: jwt.At(time.Now().Add(time.Minute * (30)))},
		user.ID,
	})
	return token.SignedString(a.secretWord)
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

func (a *UseCase) Logout(accessToken string) (string, error) {
	UserID, err := parseToken(accessToken, a.secretWord)
	if err != nil {
		return "", err
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		jwt.StandardClaims{ExpiresAt: jwt.At(time.Now().Add(time.Minute * (-30)))},
		UserID,
	})
	return expiredToken.SignedString(a.secretWord)
}

func (a *UseCase) GetCSRFToken(cookie string, expirationTime int64) (string, error) {
	return csrf.Token.Create(cookie,expirationTime)
}