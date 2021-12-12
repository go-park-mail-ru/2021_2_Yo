package usecase

import (
	"backend/pkg/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) SignUp(u *models.User) (string, error) {
	args := m.Called(u)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) SignIn(u *models.User) (string, error) {
	args := m.Called(u)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) CreateSession(userId string) (string, error) {
	args := m.Called(userId)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) CheckSession(SessionId string) (string, error) {
	args := m.Called(SessionId)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) DeleteSession(SessionId string) error {
	args := m.Called(SessionId)
	return args.Error(0)
}

func (m *UseCaseMock) CreateToken(userId string) (string, error) {
	args := m.Called(userId)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) CheckToken(csrfToken string) (string, error) {
	args := m.Called(csrfToken)
	return args.Get(0).(string), args.Error(1)
}
