package mock

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (s *RepositoryMock) CreateUser(user *models.User) error {
	args := s.Called(user)
	return args.Error(0)
}

func (s *RepositoryMock) GetUser(mail, password string) (*models.User, error) {
	args := s.Called(mail, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (s *RepositoryMock) GetUserById(userId string) (*models.User, error) {
	args := s.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}
