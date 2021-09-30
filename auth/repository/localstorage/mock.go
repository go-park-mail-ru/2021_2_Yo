package localstorage

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryUserMock struct {
	mock.Mock
}

func (s *RepositoryUserMock) CreateUser(user *models.User) error {
	args := s.Called(user)
	return args.Error(0)
}

func (s *RepositoryUserMock) GetUser(mail, password string) (*models.User, error) {
	args := s.Called(mail, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (s *RepositoryUserMock) GetUserById(userId string) (*models.User, error) {
	args := s.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}
