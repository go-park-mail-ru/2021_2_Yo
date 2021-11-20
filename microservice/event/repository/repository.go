package eventRepository

import (
	log "backend/logger"
	proto "backend/microservice/event/proto"
	"backend/models"
	error2 "backend/service/event/error"
	"context"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
)

const logMessage = "service:event:repository:postgres:"

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

const (
	checkAuthorQuery         = `select author_id from "event" where id = $1`
	listQuery                = `select * from "event"`
	getEventQuery            = `select * from "event" where id = $1`
	getEventsFromAuthorQuery = `select * from "event" where author_id = $1`
	createEventQuery         = `insert into "event" 
		(title, description, text, city, category, viewed, img_url, date, geo, tag, author_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::varchar[], $11) 
		returning id`
	updateEventQuery = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, 
		viewed = $6, img_url = $7, date = $8, geo = $9, tag = $10 
		where event.id = $11`
	updateEventQueryWithoutImgUrl = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, 
		viewed = $6, date = $7, geo = $8, tag = $9 
		where event.id = $10`
	deleteEventQuery = `delete from "event" where id = $1`
)

/*
type RepositoryServer interface {
	CreateEvent(context.Context, *Event) (*proto.EventId, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*proto.Empty, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*proto.Empty, error)
	GetEventById(context.Context, *EventId) (*proto.Event, error)
	GetEvents(context.Context, *GetEventsRequest) (*proto.Events, error)
	GetEventsFromAuthor(context.Context, *AuthorId) (*proto.Events, error)
}
*/

func (s *Repository) checkAuthor(eventId int, userId int) error {
	var authorId int
	query := checkAuthorQuery
	err := s.db.Get(&authorId, query, eventId)
	if err != nil {
		return error2.ErrPostgres
	}
	if authorId == userId {
		return nil
	} else {
		return error2.ErrNotAllowed
	}
}

func (s *Repository) CreateEvent(ctx context.Context, in *proto.Event) (*proto.EventId, error) {

	e := fromProtoToModel(in)

	newEvent, err := toPostgresEvent(e)
	if err != nil {
		return nil, err
	}

	var eventId int
	query := createEventQuery
	err = s.db.Get(&eventId, query,
		newEvent.Title,
		newEvent.Description,
		newEvent.Text,
		newEvent.City,
		newEvent.Category,
		newEvent.Viewed,
		newEvent.ImgUrl,
		newEvent.Date,
		newEvent.Geo,
		newEvent.Tag,
		newEvent.AuthorID)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrNoRows
		}
		return nil, error2.ErrPostgres
	}

	out := &proto.EventId{ID: strconv.Itoa(eventId)}
	return out, nil

}

func (s *Repository) UpdateEvent(ctx context.Context, in *proto.UpdateEventRequest) (*proto.Empty, error) {

	e := fromProtoToModel(in.Event)
	userId := in.UserId

	eventIdInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return nil, err
	}
	postgresEvent, err := toPostgresEvent(e)
	if err != nil {
		return nil, err
	}
	postgresEvent.ID = eventIdInt

	var query string
	if postgresEvent.ImgUrl != "" {
		query = updateEventQuery
		_, err = s.db.Query(query,
			postgresEvent.Title,
			postgresEvent.Description,
			postgresEvent.Text,
			postgresEvent.City,
			postgresEvent.Category,
			postgresEvent.Viewed,
			postgresEvent.ImgUrl,
			postgresEvent.Date,
			postgresEvent.Geo,
			postgresEvent.Tag,
			postgresEvent.ID)
		if err != nil {
			return nil, error2.ErrPostgres
		}
	} else {
		query = updateEventQueryWithoutImgUrl
		_, err = s.db.Query(query,
			postgresEvent.Title,
			postgresEvent.Description,
			postgresEvent.Text,
			postgresEvent.City,
			postgresEvent.Category,
			postgresEvent.Viewed,
			postgresEvent.Date,
			postgresEvent.Geo,
			postgresEvent.Tag,
			postgresEvent.ID)
		if err != nil {
			return nil, error2.ErrPostgres
		}
	}

	return nil, nil

}

func (s *Repository) DeleteEvent(ctx context.Context, in *proto.DeleteEventRequest) (*proto.Empty, error) {

	eventId := in.EventId
	userId := in.UserId

	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return nil, err
	}

	query := deleteEventQuery
	_, err = s.db.Exec(query, eventIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}

	return nil, nil

}

func (s *Repository) GetEventById(ctx context.Context, in *proto.EventId) (*proto.Event, error) {

	eventId := in.ID

	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	query := getEventQuery
	var e Event
	err = s.db.Get(&e, query, eventIdInt)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrNoRows
		}
		return nil, error2.ErrPostgres
	}
	resultEvent := toModelEvent(&e)

	out := toProtoEvent(resultEvent)
	return out, nil

}

func (s *Repository) GetEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.Events, error) {

	tags := in.Tags
	title := in.Title
	category := in.Category

	postgresTags := make(pq.StringArray, len(tags))
	for i := range tags {
		postgresTags[i] = tags[i]
	}

	query := listQuery + " "
	if title != "" {
		query += `where lower(title) ~ lower($1) and `
	} else {
		query += `where $1 = $1 and `
	}
	if category != "" {
		query += `lower(category) = lower($2) and `
	} else {
		query += `$2 = $2 and `
	}
	if len(postgresTags) != 0 {
		query += `tag && $3::varchar[]`
	} else {
		query += `$3 = $3`
	}
	query += "order by viewed DESC"
	rows, err := s.db.Queryx(query, title, category, postgresTags)
	if err != nil {
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}

	outEvents := make([]*proto.Event, len(resultEvents))
	for i, event := range resultEvents {
		outEvents[i] = toProtoEvent(event)
	}
	out := &proto.Events{Event: outEvents}
	return out, nil
}

func (s *Repository) GetEventsFromAuthor(ctx context.Context, in *proto.AuthorId) (*proto.Events, error) {

	authorId := in.ID

	authorIdInt, err := strconv.Atoi(authorId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	query := getEventsFromAuthorQuery
	rows, err := s.db.Queryx(query, authorIdInt)
	if err != nil {
		log.Error(logMessage+"GetEvents:err =", err)
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}

	outEvents := make([]*proto.Event, len(resultEvents))
	for i, event := range resultEvents {
		outEvents[i] = toProtoEvent(event)
	}
	out := &proto.Events{Event: outEvents}
	return out, nil
}