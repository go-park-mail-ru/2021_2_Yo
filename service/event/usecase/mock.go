package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) CreateEvent(e *models.Event, userId string) (string, error) {
	args := m.Called(e, userId)
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

func (m *UseCaseMock) GetEvents(title string, category string, tags []string) ([]*models.Event, error) {
	args := m.Called(title, category, tags)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *UseCaseMock) GetEventsFromAuthor(authorId string) ([]*models.Event, error) {
	args := m.Called(authorId)
	return args.Get(0).([]*models.Event), args.Error(1)
}