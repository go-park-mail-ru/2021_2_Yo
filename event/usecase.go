package event

import "backend/models"

type UseCase interface {
	List() ([]*models.Event, error)
}
