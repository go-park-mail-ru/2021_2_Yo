package event

import (
	"backend/models"
)

type Repository interface {
	GetEvent(eventId string) (*models.Event, error)
	GetEvents(title string, category string, tags []string) ([]*models.Event, error)
	CreateEvent(e *models.Event) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
}
