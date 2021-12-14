package usecase

import (
	"backend/internal/models"
	"backend/internal/service/user"
	error2 "backend/internal/service/user/error"
	"backend/internal/utils"
)

const logMessage = "service:user:usecase:"

type UseCase struct {
	repository user.Repository
}

func NewUseCase(repository user.Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (a *UseCase) GetUserById(userId string) (*models.User, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	resultUser, err := a.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	resultUser.Password = ""
	return resultUser, nil
}

func (a *UseCase) UpdateUserInfo(u *models.User) error {
	if u.ID == "" || u.Name == "" || u.Surname == "" {
		return error2.ErrEmptyData
	}
	return a.repository.UpdateUserInfo(u)
}

func (a *UseCase) UpdateUserPassword(userId string, password string) error {
	if userId == "" || password == "" {
		return error2.ErrEmptyData
	}
	hashedPassword := utils.CreatePasswordHash(password)
	return a.repository.UpdateUserPassword(userId, hashedPassword)
}

func (a *UseCase) GetSubscribers(userId string) ([]*models.User, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	resultUsers, err := a.repository.GetSubscribers(userId)
	if err != nil {
		return nil, err
	}
	for i, _ := range resultUsers {
		resultUsers[i].Password = ""
	}
	return resultUsers, nil
}

func (a *UseCase) GetSubscribes(userId string) ([]*models.User, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	resultUsers, err := a.repository.GetSubscribes(userId)
	if err != nil {
		return nil, err
	}
	for i, _ := range resultUsers {
		resultUsers[i].Password = ""
	}
	return resultUsers, nil
}

func (a *UseCase) GetFriends(userId string) ([]*models.User, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	resultUsers, err := a.repository.GetFriends(userId)
	if err != nil {
		return nil, err
	}
	for i, _ := range resultUsers {
		resultUsers[i].Password = ""
	}
	return resultUsers, nil
}

func (a *UseCase) GetVisitors(eventId string) ([]*models.User, error) {
	if eventId == "" {
		return nil, error2.ErrEmptyData
	}
	resultUsers, err := a.repository.GetVisitors(eventId)
	if err != nil {
		return nil, err
	}
	for i, _ := range resultUsers {
		resultUsers[i].Password = ""
	}
	return resultUsers, nil
}

func (a *UseCase) Subscribe(subscribedId string, subscriberId string) error {
	if subscribedId == "" || subscriberId == "" {
		return error2.ErrEmptyData
	}
	return a.repository.Subscribe(subscribedId, subscriberId)
}

func (a *UseCase) Unsubscribe(subscribedId string, subscriberId string) error {
	if subscribedId == "" || subscriberId == "" {
		return error2.ErrEmptyData
	}
	return a.repository.Unsubscribe(subscribedId, subscriberId)
}

func (a *UseCase) IsSubscribed(subscribedId string, subscriberId string) (bool, error) {
	if subscribedId == "" || subscriberId == "" {
		return false, error2.ErrEmptyData
	}
	if subscribedId == subscriberId {
		return false, nil
	}
	return a.repository.IsSubscribed(subscribedId, subscriberId)
}
