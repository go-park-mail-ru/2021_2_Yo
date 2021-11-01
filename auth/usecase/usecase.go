package usecase

import (
	"backend/auth"
	error2 "backend/auth/error"
	"backend/models"
	"crypto/sha256"
	"fmt"
)

const logMessage = "auth:usecase:usecase:"

func createPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

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

func (a *UseCase) SignUp(user *models.User) (string, error) {
	if user == nil {
		return "", error2.ErrEmptyData
	}
	hashedPassword := createPasswordHash(user.Password)
	user.Password = hashedPassword
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail string, password string) (string, error) {
	if mail == "" || password == "" {
		return "", error2.ErrEmptyData
	}
	hashedPassword := createPasswordHash(password)
	user, err := a.repository.GetUser(mail, hashedPassword)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (a *UseCase) GetUser(userId string) (*models.User, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.repository.GetUserById(userId)
}

func (a *UseCase) UpdateUserInfo(userId string, name string, surname string, about string) error {
	if userId == "" {
		return error2.ErrEmptyData
	}
	return a.repository.UpdateUserInfo(userId, name, surname, about)
}

func (a *UseCase) UpdateUserPassword(userId string, password string) error {
	if userId == "" || password == "" {
		return error2.ErrEmptyData
	}
	hashedPassword := createPasswordHash(password)
	return a.repository.UpdateUserPassword(userId, hashedPassword)
}
