package mock

import (
	"backend/internal/models"
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

func (m *RepositoryMock) GetEvents(userId string, title string, category string, city string, date string, tags []string) ([]*models.Event, error) {
	args := m.Called(userId, title, category, city, date, tags)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *RepositoryMock) GetCreatedEvents(authorId string) ([]*models.Event, error) {
	args := m.Called(authorId)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *RepositoryMock) GetVisitedEvents(userId string) ([]*models.Event, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *RepositoryMock) Visit(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *RepositoryMock) Unvisit(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *RepositoryMock) IsVisited(eventId string, userId string) (bool, error) {
	args := m.Called(eventId, userId)
	return args.Get(0).(bool), args.Error(1)
}

func (m *RepositoryMock) GetCities() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *RepositoryMock) EmailNotify(eventId string) ([]*models.Info, error) {
	args := m.Called(eventId)
	return args.Get(0).([]*models.Info), args.Error(1)
}
