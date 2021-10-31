package usecase

import (
	"backend/event"
	log "backend/logger"
	"backend/models"
	"errors"
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
		err := errors.New("eventId is null")
		return nil, err
	}
	return a.eventRepo.GetEvent(eventId)
}

func (a *UseCase) CreateEvent(event *models.Event, userId string) (string, error) {
	if event == nil || userId == "" {
		err := errors.New("UseCase:CreateEvent event or userId is empty")
		return "", err
	}
	event.AuthorId = userId
	log.Debug("UseCase CreateEvent: HERE")
	return a.eventRepo.CreateEvent(event)
}

func (a *UseCase) UpdateEvent(event *models.Event, userId string) error {
	if event == nil || userId == "" {
		err := errors.New("event or userId is nil")
		return err
	}
	if event.ID == "" {
		err := errors.New("event.ID is null")
		return err
	}
	return a.eventRepo.UpdateEvent(event, userId)
}

func (a *UseCase) DeleteEvent(userId string, eventID string) error {
	if userId == "" || eventID == "" {
		err := errors.New("eventId is empty")
		return err
	}
	return a.eventRepo.DeleteEvent(userId, eventID)
}
