package eventRepository

import (
	proto "backend/microservice/event/proto"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	"strconv"

	"github.com/lib/pq"
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
	Address		string         `db:"address"`
	AuthorID    int            `db:"author_id"`
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
		Address: 	 e.Address,
		AuthorID:    authorIdInt,
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
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		Address: 	 e.Address,
		AuthorId:    strconv.Itoa(e.AuthorID),
	}
}

func toProtoEvent(e *models.Event) *proto.Event {
	return &proto.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      int32(e.Viewed),
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		Address: 	 e.Address,
		AuthorId:    e.AuthorId,
	}
}

func fromProtoToModel(in *proto.Event) *models.Event {
	return &models.Event{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Text:        in.Text,
		City:        in.City,
		Category:    in.Category,
		Viewed:      int(in.Viewed),
		ImgUrl:      in.ImgUrl,
		Tag:         in.Tag,
		Date:        in.Date,
		Geo:         in.Geo,
		Address: 	 in.Address,
		AuthorId:    in.AuthorId,
	}
}
