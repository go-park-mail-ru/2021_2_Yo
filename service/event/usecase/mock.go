package usecase

import (
	"backend/pkg/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) CreateEvent(e *models.Event) (string, error) {
	args := m.Called(e)
	return args.Get(0).(string), args.Error(1)
}

func (m *UseCaseMock) UpdateEvent(e *models.Event, userId string) error {
	args := m.Called(e, userId)
	return args.Error(0)
}

func (m *UseCaseMock) DeleteEvent(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *UseCaseMock) GetEventById(eventId string) (*models.Event, error) {
	args := m.Called(eventId)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *UseCaseMock) GetEvents(title string, category string, city string, date string, tags []string) ([]*models.Event, error) {
	args := m.Called(title, category, city, date, tags)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *UseCaseMock) GetCreatedEvents(authorId string) ([]*models.Event, error) {
	args := m.Called(authorId)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *UseCaseMock) GetVisitedEvents(userId string) ([]*models.Event, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *UseCaseMock) Visit(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *UseCaseMock) Unvisit(eventId string, userId string) error {
	args := m.Called(eventId, userId)
	return args.Error(0)
}

func (m *UseCaseMock) IsVisited(eventId string, userId string) (bool, error) {
	args := m.Called(eventId, userId)
	return args.Get(0).(bool), args.Error(1)
}

func (m *UseCaseMock) GetCities() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
