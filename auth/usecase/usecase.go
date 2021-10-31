package usecase

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
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

func (a *UseCase) SignUp(user *models.User) (string, error) {
	hashedPassword := a.CreatePasswordHash(user.Password)
	user.Password = hashedPassword
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail, password string) (string, error) {
	message := logMessage + "SignIn:"
	hashedPassword := a.CreatePasswordHash(password)
	user, err := a.repository.GetUser(mail, hashedPassword)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		return "", err
	}
	return user.ID, nil
}

func (a *UseCase) GetUser(userId string) (*models.User, error) {
	return a.repository.GetUserById(userId)
}

func (a *UseCase) CreatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (a *UseCase) Logout(accessToken string) (string, error) {
	/*
		UserID, err := parseToken(accessToken, a.secretWord)
		if err != nil {
			return "", err
		}
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
			jwt.StandardClaims{ExpiresAt: jwt.At(time.Now().Add(time.Minute * (-30)))},
			UserID,
		})
		return expiredToken.SignedString(a.secretWord)
	*/
	return "", nil
}
