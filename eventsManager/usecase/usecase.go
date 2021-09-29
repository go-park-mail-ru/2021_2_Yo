package usecase

import (
	"backend/eventsManager"
	"backend/models"
)

type UseCaseEventsManager struct {
	eventRepo eventsManager.RepositoryEventsManager
}

func NewUseCaseEvents(eventRepo eventsManager.RepositoryEventsManager) *UseCaseEventsManager {
	return &UseCaseEventsManager{
		eventRepo: eventRepo,
	}
}

func (a *UseCaseEventsManager) List() ([]*models.Event, error) {
	return a.eventRepo.List()
}
