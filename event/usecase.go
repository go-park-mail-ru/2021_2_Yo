package event

import "backend/models"

type UseCase interface {
	List() ([]*models.Event, error)
	GetEvent(string) (*models.Event, error)
	UpdateEvent(string, *models.Event) error
	CreateEvent(*models.Event) (string, error)
}
