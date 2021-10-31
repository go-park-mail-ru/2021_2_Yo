package usecase

import (
	"backend/auth"
	"backend/models"
	"crypto/sha256"
	"errors"
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
		err := errors.New("user is nil")
		return "", err
	}
	hashedPassword := createPasswordHash(user.Password)
	user.Password = hashedPassword
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail string, password string) (string, error) {
	if mail == "" || password == "" {
		err := errors.New("mail or password is nil")
		return "", err
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
		err := errors.New("userId is empty")
		return nil, err
	}
	return a.repository.GetUserById(userId)
}

func (a *UseCase) UpdateUserInfo(userId string, name string, surname string, about string) error {
	if userId == "" || name == "" || surname == "" || about == "" {
		err := errors.New("UpdateUserInfo data in empty")
		return err
	}
	return a.repository.UpdateUserInfo(userId, name, surname, about)
}

func (a *UseCase) UpdateUserPassword(userId string, password string) error {
	if userId == "" || password == "" {
		err := errors.New("UpdateUserPassword data in empty")
		return err
	}
	hashedPassword := createPasswordHash(password)
	return a.repository.UpdateUserPassword(userId, hashedPassword)
}
