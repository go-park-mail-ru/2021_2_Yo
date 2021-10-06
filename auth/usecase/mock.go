package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) SignUp(name, surname, mail, password string) error {
	args := m.Called(name, surname, mail, password)
	return args.Error(0)
}

func (m *UseCaseMock) SignIn(mail, password string) (string, error) {
	args := m.Called(mail, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) ParseToken(accessToken string) (*models.User, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*models.User), args.Error(1)
}
