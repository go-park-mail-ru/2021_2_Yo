package usecase

import (
	"backend/event"
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
	return a.eventRepo.GetEvent(eventId)
}

func (a *UseCase) UpdateEvent(eventID string, event *models.Event) error {
	return a.eventRepo.UpdateEvent(eventID, event)
}

func (a *UseCase) CreateEvent(event *models.Event) (string, error) {
	return a.eventRepo.CreateEvent(event)
}
