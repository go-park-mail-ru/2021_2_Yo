package eventsManager

import (
	"backend/models"
)

type RepositoryEventsManager interface {
	List() ([]*models.Event, error)
}
