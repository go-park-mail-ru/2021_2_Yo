package usecase

import (
	"backend/models"
	"backend/service/event"
	error2 "backend/service/event/error"
	"strings"
)

const logMessage = "service:event:usecase:"

type UseCase struct {
	eventRepo event.Repository
}

func NewUseCase(eventRepo event.Repository) *UseCase {
	return &UseCase{
		eventRepo: eventRepo,
	}
}

func (a *UseCase) CreateEvent(e *models.Event) (string, error) {
	if e == nil || e.AuthorId == "" {
		return "", error2.ErrEmptyData
	}
	for _, tag := range e.Tag {
		tag = strings.ToLower(tag)
	}
	return a.eventRepo.CreateEvent(e)
}

func (a *UseCase) UpdateEvent(e *models.Event, userId string) error {
	if e == nil || userId == "" {
		return error2.ErrEmptyData
	}
	if e.ID == "" {
		return error2.ErrEmptyData
	}
	for _, tag := range e.Tag {
		tag = strings.ToLower(tag)
	}
	return a.eventRepo.UpdateEvent(e, userId)
}

func (a *UseCase) DeleteEvent(userId string, eventID string) error {
	if userId == "" || eventID == "" {
		return error2.ErrEmptyData
	}
	return a.eventRepo.DeleteEvent(userId, eventID)
}

func (a *UseCase) GetEventById(eventId string) (*models.Event, error) {
	if eventId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.eventRepo.GetEventById(eventId)
}

func (a *UseCase) GetEvents(title string, category string, tags []string) ([]*models.Event, error) {
	if tags != nil && tags[0] == "" {
		tags = nil
	}
	return a.eventRepo.GetEvents(title, category, tags)
}

func (a *UseCase) GetEventsFromAuthor(authorId string) ([]*models.Event, error) {
	if authorId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.eventRepo.GetEventsFromAuthor(authorId)
}
