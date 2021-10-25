package event

import (
	"backend/models"
)

type Repository interface {
	List() ([]*models.Event, error)
	GetEvent(eventId string) (*models.Event, error)
	CreateEvent(event *models.Event) (string, error)
}
