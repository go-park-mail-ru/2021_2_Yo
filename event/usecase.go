package event

import "backend/models"

type UseCase interface {
	List() ([]*models.Event, error)
	GetEvent(string) (*models.Event, error)
	CreateEvent(*models.Event, string) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
}
