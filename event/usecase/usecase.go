package usecase

import (
	"backend/event"
	error2 "backend/event/error"
	"backend/models"
)

type UseCase struct {
	eventRepo event.Repository
}

func NewUseCase(eventRepo event.Repository) *UseCase {
	return &UseCase{
		eventRepo: eventRepo,
	}
}

func (a *UseCase) List() ([]*models.Event, error) {
	return a.eventRepo.List()
}

func (a *UseCase) GetEvent(eventId string) (*models.Event, error) {
	if eventId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.eventRepo.GetEvent(eventId)
}

func (a *UseCase) CreateEvent(e *models.Event, userId string) (string, error) {
	if e == nil || userId == "" {
		return "", error2.ErrEmptyData
	}
	e.AuthorId = userId
	return a.eventRepo.CreateEvent(e)
}

func (a *UseCase) UpdateEvent(e *models.Event, userId string) error {
	if e == nil || userId == "" {
		return error2.ErrEmptyData
	}
	if e.ID == "" {
		return error2.ErrEmptyData
	}
	return a.eventRepo.UpdateEvent(e, userId)
}

func (a *UseCase) DeleteEvent(userId string, eventID string) error {
	if userId == "" || eventID == "" {
		return error2.ErrEmptyData
	}
	return a.eventRepo.DeleteEvent(userId, eventID)
}
