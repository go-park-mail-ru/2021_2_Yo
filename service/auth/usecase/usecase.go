package usecase

import (
	"backend/models"
	"backend/service/auth"
	error2 "backend/service/auth/error"
	"backend/utils"
)

const logMessage = "service:auth:usecase:"

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
	hashedPassword := utils.CreatePasswordHash(user.Password)
	user.Password = hashedPassword
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail string, password string) (string, error) {
	if mail == "" || password == "" {
		return "", error2.ErrEmptyData
	}
	hashedPassword := utils.CreatePasswordHash(password)
	user, err := a.repository.GetUser(mail, hashedPassword)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}
