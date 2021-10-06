package localstorage

import (
	"backend/models"
	"strconv"
)

type Event struct {
	ID          int
	Name        string
	Description string
	Views       int
	ImgUrl      string
}

func toLocalstorageEvent(e *models.Event) *Event {
	return &Event{
		Name:        e.Name,
		Description: e.Description,
		Views:       e.Views,
		ImgUrl:      e.ImgUrl,
	}
}

func toModelEvent(e *Event) *models.Event {
	return &models.Event{
		ID:          strconv.Itoa(e.ID),
		Name:        e.Name,
		Description: e.Description,
		Views:       e.Views,
		ImgUrl:      e.ImgUrl,
	}
}
