package postgres

import (
	"backend/event"
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
	Img_Url     string
	Tag         []string
	Date        string
	Geo         string
	Author_ID   int
}

func toPostgresEvent(e *models.Event) (*Event, error) {
	authorIdInt := 0
	if e.AuthorId == "" {
		authorIdInt = 0
	} else {
		authorIdInt, err := strconv.Atoi(e.AuthorId)
		if err != nil {
			return nil, event.ErrAtoi
		}
		_ = authorIdInt
	}
	return &Event{
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      e.Viewed,
		Img_Url:     e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		Author_ID:   authorIdInt,
	}, nil
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
		ImgUrl:      e.Img_Url,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		AuthorId:    strconv.Itoa(e.Author_ID),
	}
}
