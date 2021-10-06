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
