package usecase

import (
	"backend/models"
	user "backend/service/user"
	error2 "backend/service/user/error"
	"backend/utils"
)

const logMessage = "service:user:usecase:"

type UseCase struct {
	repository user.Repository
}

func NewUseCase(userRepo user.Repository) *UseCase {
	return &UseCase{
		repository: userRepo,
	}
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
	hashedPassword := utils.CreatePasswordHash(password)
	return a.repository.UpdateUserPassword(userId, hashedPassword)
}
