package postgres

import (
	error2 "backend/internal/service/event/error"
	"backend/pkg/models"
	"github.com/lib/pq"
	"strconv"
)

type Event struct {
	ID          int            `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Text        string         `db:"text"`
	City        string         `db:"city"`
	Category    string         `db:"category"`
	Viewed      int            `db:"viewed"`
	ImgUrl      string         `db:"img_url"`
	Tag         pq.StringArray `db:"tag"`
	Date        string         `db:"date"`
	Geo         string         `db:"geo"`
	Address     string         `db:"address"`
	AuthorID    int            `db:"author_id"`
	IsVisited   int            `db:"count"`
}

func toPostgresEvent(e *models.Event) (*Event, error) {
	var authorIdInt int
	if e.AuthorId == "" {
		authorIdInt = 0
	} else {
		tempAuthorId, err := strconv.Atoi(e.AuthorId)
		if err != nil {
			return nil, error2.ErrAtoi
		}
		authorIdInt = tempAuthorId
	}
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
		Address:     e.Address,
		AuthorID:    authorIdInt,
	}, nil
}

func toModelEvent(e *Event) *models.Event {
	isVisited := false
	if e.IsVisited > 0 {
		isVisited = true
	}
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
		Address:     e.Address,
		AuthorId:    strconv.Itoa(e.AuthorID),
		IsVisited:   isVisited,
	}
}
