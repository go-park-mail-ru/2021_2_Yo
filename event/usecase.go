package event

import "backend/models"

type UseCase interface {
	CreateEvent(*models.Event, string) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
	GetEvent(string) (*models.Event, error)
	GetEvents(title string, category string, tags []string) ([]*models.Event, error)
	GetEventsFromAuthor(authorId string) ([]*models.Event, error)
}
