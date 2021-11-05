package mock

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (s *RepositoryMock) CreateUser(user *models.User) (string, error) {
	args := s.Called(user)
	return args.Get(0).(string), args.Error(1)
}

func (s *RepositoryMock) GetUser(mail, password string) (*models.User, error) {
	args := s.Called(mail, password)
	return args.Get(0).(*models.User), args.Error(1)
}
