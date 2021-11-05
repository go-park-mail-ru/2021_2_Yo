package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) SignUp(user *models.User) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) SignIn(mail string, password string) (string, error) {
	args := m.Called(mail, password)
	return args.Get(0).(string), args.Error(1)
}
