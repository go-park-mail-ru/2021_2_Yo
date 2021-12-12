package event

import (
	"backend/internal/models"
)

type UseCase interface {
	CreateEvent(e *models.Event) (string, error)
	UpdateEvent(e *models.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
	//
	GetEventById(eventId string) (*models.Event, error)
	GetEvents(userId string, title string, category string, city string, date string, tags []string) ([]*models.Event, error)
	GetCreatedEvents(authorId string) ([]*models.Event, error)
	GetVisitedEvents(userId string) ([]*models.Event, error)
	//
	Visit(eventId string, userId string) error
	Unvisit(eventId string, userId string) error
	IsVisited(eventId string, userId string) (bool, error)
	//
	GetCities() ([]string, error)
	//
	EmailNotify(eventId string) error
}
