package usecase

import (
	"backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) GetUserById(userId string) (*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UseCaseMock) UpdateUserInfo(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UseCaseMock) UpdateUserPassword(userId string, password string) error {
	args := m.Called(userId, password)
	return args.Error(0)
}

func (m *UseCaseMock) GetSubscribers(userId string) ([]*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UseCaseMock) GetSubscribes(userId string) ([]*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UseCaseMock) GetFriends(userId string, eventId string) ([]*models.User, error) {
	args := m.Called(userId, eventId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UseCaseMock) GetVisitors(eventId string) ([]*models.User, error) {
	args := m.Called(eventId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UseCaseMock) Subscribe(subscribedId string, subscriberId string) error {
	args := m.Called(subscribedId, subscriberId)
	return args.Error(0)
}

func (m *UseCaseMock) Unsubscribe(subscribedId string, subscriberId string) error {
	args := m.Called(subscribedId, subscriberId)
	return args.Error(0)
}

func (m *UseCaseMock) IsSubscribed(subscribedId string, subscriberId string) (bool, error) {
	args := m.Called(subscribedId, subscriberId)
	return args.Get(0).(bool), args.Error(1)
}
