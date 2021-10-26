package localstorage

import (
	"backend/models"
	"strconv"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Text        string
	City        string
	Category    string
	Viewed      int
	ImgUrl      string
	Tag         []string
	Date        string
	Geo         string
}

func toLocalstorageEvent(e *models.Event) *Event {
	return &Event{
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      e.Viewed,
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
	}
}

func toModelEvent(e *Event) *models.Event {
	return &models.Event{
		ID:          strconv.Itoa(e.ID),
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      e.Viewed,
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
	}
}
