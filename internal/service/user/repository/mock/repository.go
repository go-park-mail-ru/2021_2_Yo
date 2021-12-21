package mock

import (
	"backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) GetUserById(userId string) (*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *RepositoryMock) UpdateUserInfo(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateUserPassword(userId string, password string) error {
	args := m.Called(userId, password)
	return args.Error(0)
}

func (m *RepositoryMock) GetSubscribers(userId string) ([]*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *RepositoryMock) GetSubscribes(userId string) ([]*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *RepositoryMock) GetFriends(userId string, eventId string) ([]*models.User, error) {
	args := m.Called(userId, eventId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *RepositoryMock) GetVisitors(eventId string) ([]*models.User, error) {
	args := m.Called(eventId)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *RepositoryMock) Subscribe(subscribedId string, subscriberId string) error {
	args := m.Called(subscribedId, subscriberId)
	return args.Error(0)
}

func (m *RepositoryMock) Unsubscribe(subscribedId string, subscriberId string) error {
	args := m.Called(subscribedId, subscriberId)
	return args.Error(0)
}

func (m *RepositoryMock) IsSubscribed(subscribedId string, subscriberId string) (bool, error) {
	args := m.Called(subscribedId, subscriberId)
	return args.Get(0).(bool), args.Error(1)
}
