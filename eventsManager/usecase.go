package eventsManager

import "backend/models"

type UseCaseEventsManager interface {
	List() ([]*models.Event, error)
}
