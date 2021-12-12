package usecase

import (
	"backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type AuthRepoMock struct {
	mock.Mock
}

func (m *AuthRepoMock) CreateUser(data *models.User) (string, error) {
	args := m.Called(data)
	return args.Get(0).(string), args.Error(1)
}

func (m *AuthRepoMock) GetUser(mail, password string) (*models.User, error) {
	args := m.Called(mail, password)
	return args.Get(0).(*models.User), args.Error(1)
}
