package eventRepository

import (
	proto "backend/microservice/event/proto"
	log "backend/pkg/logger"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	"context"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

const (
	logMessage       = "microservice:event:repository:"
	checkAuthorQuery = `select author_id from "event" where id = $1`
	listQuery        = `select * from "event"`
	getEventQuery    = `select * from "event" where id = $1`
	createEventQuery = `insert into "event" 
		(title, description, text, city, category, viewed, img_url, date, geo, address, tag, author_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::varchar[], $12) 
		returning id`
	updateEventQuery = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, 
		viewed = $6, img_url = $7, date = $8, geo = $9, address = $10, tag = $11 
		where event.id = $12`
	updateEventQueryWithoutImgUrl = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, 
		viewed = $6, date = $7, geo = $8, address = $9, tag = $10 
		where event.id = $11`
	deleteEventQuery = `delete from "event" where id = $1`
	visitedQuery     = `select e.* from "event" as e join visitor as v on v.event_id = e.id where v.user_id = $1`
	createdQuery     = `select * from "event" where author_id = $1`
	visitQuery       = `insert into "visitor" (event_id, user_id) values ($1, $2)`
	unvisitQuery     = `delete from "visitor" where event_id = $1 and user_id = $2`
	isVisitedQuery   = `select count(*) from "visitor" where event_id = $1 and user_id = $2`
	getCitiesQuery   = `select distinct city from event`
)

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
	message := logMessage + "CreateEvent:"
	log.Debug(message + "started")
	e := fromProtoToModel(in)
	newEvent, err := toPostgresEvent(e)
	if err != nil {
		log.Error(message+"err 1 =", err)
		return &proto.EventId{}, err
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
		newEvent.Address,
		newEvent.Tag,
		newEvent.AuthorID)
	if err != nil {
		log.Error(message+"err 2 =", err)
		if err == sql2.ErrNoRows {
			return &proto.EventId{}, error2.ErrNoRows
		}
		return &proto.EventId{}, error2.ErrPostgres
	}
	out := &proto.EventId{ID: strconv.Itoa(eventId)}
	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) UpdateEvent(ctx context.Context, in *proto.UpdateEventRequest) (*proto.Empty, error) {
	message := logMessage + "UpdateEvent:"
	log.Debug(message + "started")
	e := fromProtoToModel(in.Event)
	userId := in.UserId
	eventIdInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return &proto.Empty{}, err
	}
	postgresEvent, err := toPostgresEvent(e)
	if err != nil {
		return &proto.Empty{}, err
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
			postgresEvent.Address,
			postgresEvent.Tag,
			postgresEvent.ID)
		if err != nil {
			return &proto.Empty{}, error2.ErrPostgres
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
			postgresEvent.Address,
			postgresEvent.Tag,
			postgresEvent.ID)
		if err != nil {
			return &proto.Empty{}, error2.ErrPostgres
		}
	}
	log.Debug(message + "ended")
	return &proto.Empty{}, nil
}

func (s *Repository) DeleteEvent(ctx context.Context, in *proto.DeleteEventRequest) (*proto.Empty, error) {
	message := logMessage + "DeleteEvent:"
	log.Debug(message + "started")
	eventId := in.EventId
	userId := in.UserId
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return &proto.Empty{}, err
	}
	query := deleteEventQuery
	_, err = s.db.Query(query, eventIdInt)
	if err != nil {
		log.Error(err)
		return &proto.Empty{}, error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return &proto.Empty{}, nil
}

func (s *Repository) GetEventById(ctx context.Context, in *proto.EventId) (*proto.Event, error) {
	message := logMessage + "GetEventById:"
	log.Debug(message + "started")
	eventId := in.ID
	query := getEventQuery
	var e Event
	err := s.db.Get(&e, query, eventId)
	if err != nil {
		log.Error(err)
		if err == sql2.ErrNoRows {
			return &proto.Event{}, error2.ErrNoRows
		}
		return &proto.Event{}, error2.ErrPostgres
	}
	resultEvent := toModelEvent(&e)
	out := toProtoEvent(resultEvent)
	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) GetEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.Events, error) {
	message := logMessage + "GetEvents:"
	log.Debug(message + "started")
	title := in.Title
	category := in.Category
	city := in.City
	date := in.Date
	tags := in.Tags
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
	if city != "" {
		query += `lower(city) = lower($3) and `
	} else {
		query += `$3 = $3 and `
	}
	if date != "" {
		query += `lower(date) = lower($4) and `
	} else {
		query += `$4 = $4 and `
	}
	if len(postgresTags) != 0 {
		query += `tag && $5::varchar[]`
	} else {
		query += `$5 = $5`
	}
	query += ` order by viewed DESC`
	rows, err := s.db.Queryx(query, title, category, city, date, postgresTags)
	if err != nil {
		log.Error(message, "err = ", err)
		return &proto.Events{}, err
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			log.Error(message, "err = ", err)
			return &proto.Events{}, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}
	outEvents := make([]*proto.Event, len(resultEvents))
	for i, event := range resultEvents {
		outEvents[i] = toProtoEvent(event)
	}
	out := &proto.Events{Events: outEvents}
	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) GetVisitedEvents(ctx context.Context, in *proto.UserId) (*proto.Events, error) {
	message := logMessage + "GetVisitedEvents:"
	log.Debug(message + "started")
	userId := in.ID
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Events{}, error2.ErrAtoi
	}
	query := visitedQuery
	rows, err := s.db.Queryx(query, userIdInt)
	if err != nil {
		return &proto.Events{}, error2.ErrPostgres
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			return &proto.Events{}, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}
	outEvents := make([]*proto.Event, len(resultEvents))
	for i, event := range resultEvents {
		outEvents[i] = toProtoEvent(event)
	}
	out := &proto.Events{Events: outEvents}
	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) GetCreatedEvents(ctx context.Context, in *proto.UserId) (*proto.Events, error) {
	message := logMessage + "GetCreatedEvents:"
	log.Debug(message + "started")
	userId := in.ID
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Events{}, error2.ErrAtoi
	}
	query := createdQuery
	rows, err := s.db.Queryx(query, userIdInt)
	if err != nil {
		return &proto.Events{}, error2.ErrPostgres
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			return &proto.Events{}, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}
	outEvents := make([]*proto.Event, len(resultEvents))
	for i, event := range resultEvents {
		outEvents[i] = toProtoEvent(event)
	}
	out := &proto.Events{Events: outEvents}
	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) Visit(ctx context.Context, in *proto.VisitRequest) (*proto.Empty, error) {
	message := logMessage + "Visit:"
	log.Debug(message + "started")
	eventId := in.EventId
	userId := in.UserId
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	query := visitQuery
	_, err = s.db.Query(query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return &proto.Empty{}, error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return &proto.Empty{}, nil
}

func (s *Repository) Unvisit(ctx context.Context, in *proto.VisitRequest) (*proto.Empty, error) {
	message := logMessage + "Unvisit:"
	log.Debug(message + "started")
	eventId := in.EventId
	userId := in.UserId
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.Empty{}, error2.ErrAtoi
	}
	query := unvisitQuery
	_, err = s.db.Query(query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return &proto.Empty{}, error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return &proto.Empty{}, nil
}

func (s *Repository) IsVisited(ctx context.Context, in *proto.VisitRequest) (*proto.IsVisitedRequest, error) {
	message := logMessage + "IsVisited:"
	log.Debug(message + "started")
	eventId := in.EventId
	userId := in.UserId
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return &proto.IsVisitedRequest{}, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &proto.IsVisitedRequest{}, error2.ErrAtoi
	}
	query := isVisitedQuery
	var count int
	result := false
	err = s.db.Get(&count, query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return &proto.IsVisitedRequest{}, error2.ErrPostgres
	}
	if count > 0 {
		result = true
	}
	log.Debug(message + "ended")
	return &proto.IsVisitedRequest{
		Result: result,
	}, nil
}

func (s *Repository) GetCities(ctx context.Context, in *proto.Empty) (*proto.GetCitiesRequest, error) {
	message := logMessage + "GetCities:"
	log.Debug(message + "started")
	query := getCitiesQuery
	rows, err := s.db.Queryx(query)
	if err != nil {
		return &proto.GetCitiesRequest{}, error2.ErrPostgres
	}
	defer rows.Close()
	var resultCities []string
	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			return &proto.GetCitiesRequest{}, error2.ErrPostgres
		}
		resultCities = append(resultCities, c)
	}
	log.Debug(message + "ended")
	return &proto.GetCitiesRequest{
		Cities: resultCities,
	}, nil
}
