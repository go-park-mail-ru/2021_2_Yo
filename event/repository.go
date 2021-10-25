package event

import (
	"backend/models"
)

type Repository interface {
	List() ([]*models.Event, error)
	GetEvent(eventId string) (*models.Event, error)
}
