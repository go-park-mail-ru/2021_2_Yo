package event

import (
	"backend/models"
)

type Repository interface {
	List() ([]*models.Event, error)
}
