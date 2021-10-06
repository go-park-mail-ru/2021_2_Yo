package mock

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (s *RepositoryMock) List() ([]*models.Event, error) {
	return nil, nil
}
