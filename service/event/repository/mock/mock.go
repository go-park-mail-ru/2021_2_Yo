package mock

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) CreateEvent(e *models.Event) (string, error) {
	args := m.Called(e)
	return args.Get(0).(string), args.Error(1)
}

func (m *RepositoryMock) UpdateEvent(e *models.Event, userId string) error {
	args := m.Called(e, userId)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteEvent(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *RepositoryMock) GetEventById(eventId string) (*models.Event, error) {
	args := m.Called(eventId)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *RepositoryMock) GetEvents(title string, category string, tags []string) ([]*models.Event, error) {
	args := m.Called(title, category, tags)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *RepositoryMock) GetEventsFromAuthor(authorId string) ([]*models.Event, error) {
	args := m.Called(authorId)
	return args.Get(0).([]*models.Event), args.Error(1)
}
