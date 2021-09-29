package localstorage

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryEventMock struct {
	mock.Mock
}

func (s *RepositoryEventMock) List() ([]*models.Event, error) {
	return nil, nil
}
