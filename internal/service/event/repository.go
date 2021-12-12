package event

import (
	models2 "backend/internal/models"
)

type Repository interface {
	CreateEvent(e *models2.Event) (string, error)
	UpdateEvent(e *models2.Event, userId string) error
	DeleteEvent(eventId string, userId string) error
	//
	GetEventById(eventId string) (*models2.Event, error)
	GetEvents(userId string, title string, category string, city string, date string, tags []string) ([]*models2.Event, error)
	GetCreatedEvents(authorId string) ([]*models2.Event, error)
	GetVisitedEvents(userId string) ([]*models2.Event, error)
	//
	Visit(eventId string, userId string) error
	Unvisit(eventId string, userId string) error
	IsVisited(eventId string, userId string) (bool, error)
	//
	GetCities() ([]string, error)
	//
	EmailNotify(eventId string) ([]*models2.Info, error)
}
