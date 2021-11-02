package mock

import (
	error2 "backend/event/error"
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (s *RepositoryMock) List() ([]*models.Event, error) {
	return nil, nil
}

func (s *RepositoryMock) GetEvent(eventId string) (*models.Event, error) {
	return nil, error2.ErrEventNotFound
}
