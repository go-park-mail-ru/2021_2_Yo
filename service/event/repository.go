package event

import (
	"backend/models"
)

type Repository interface {
	CreateEvent(e *models.Event) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
	GetEventById(eventId string) (*models.Event, error)
	GetEvents(title string, category string, tags []string) ([]*models.Event, error)
	GetEventsFromAuthor(authorId string) ([]*models.Event, error)
}
