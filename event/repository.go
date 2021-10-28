package event

import (
	"backend/models"
)

type Repository interface {
	List() ([]*models.Event, error)
	GetEvent(eventId string) (*models.Event, error)
	CreateEvent(event *models.Event) (string, error)
	UpdateEvent(eventId string, e *models.Event) error
	DeleteEvent(eventId string, userId string) error
}
