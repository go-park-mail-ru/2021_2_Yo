package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) GetUser(userId string) (*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UseCaseMock) UpdateUserInfo(userId string, name string, surname string, about string) error {
	args := m.Called(userId, name, surname, about)
	return args.Error(0)
}

func (m *UseCaseMock) UpdateUserPassword(userId string, password string) error {
	args := m.Called(userId, password)
	return args.Error(0)
}
