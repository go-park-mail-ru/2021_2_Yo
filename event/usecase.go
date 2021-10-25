package event

import "backend/models"

type UseCase interface {
	List() ([]*models.Event, error)
	Event(eventId string) (*models.Event, error)
}
