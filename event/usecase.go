package event

import "backend/models"

type UseCase interface {
	GetEvent(string) (*models.Event, error)
	GetEvents(title string, category string, tags []string) ([]*models.Event, error)
	CreateEvent(*models.Event, string) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
}
