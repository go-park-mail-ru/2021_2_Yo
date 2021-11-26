package postgres

import (
	log "backend/pkg/logger"
	"backend/pkg/models"
	error2 "backend/service/event/error"
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
	logMessage       = "service:event:repository:postgres:"
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

func (s *Repository) CreateEvent(e *models.Event) (string, error) {
	message := logMessage + "CreateEvent:"
	log.Debug(message + "started")
	newEvent, err := toPostgresEvent(e)
	if err != nil {
		return "", err
	}
	var eventId int
	query := createEventQuery
	err = s.db.Get(
		&eventId, query,
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
		if err == sql2.ErrNoRows {
			return "", error2.ErrNoRows
		}
		return "", error2.ErrPostgres
	}
	eventIdStr := strconv.Itoa(eventId)
	log.Debug(message + "ended")
	return eventIdStr, nil
}

func (s *Repository) UpdateEvent(e *models.Event, userId string) error {
	message := logMessage + "UpdateEvent:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return err
	}
	postgresEvent, err := toPostgresEvent(e)
	if err != nil {
		return err
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
			return error2.ErrPostgres
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
			return error2.ErrPostgres
		}
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) DeleteEvent(eventId string, userId string) error {
	message := logMessage + "DeleteEvent:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	err = s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		return err
	}
	query := deleteEventQuery
	_, err = s.db.Query(query, eventIdInt)
	if err != nil {
		log.Error(err)
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) GetEventById(eventId string) (*models.Event, error) {
	message := logMessage + "GetEventById:"
	log.Debug(message + "started")
	query := getEventQuery
	var e Event
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	err = s.db.Get(&e, query, eventIdInt)
	if err != nil {
		log.Error(err)
		if err == sql2.ErrNoRows {
			return nil, error2.ErrNoRows
		}
		return nil, error2.ErrPostgres
	}
	resultEvent := toModelEvent(&e)
	log.Debug(message + "ended")
	return resultEvent, nil
}

func (s *Repository) GetEvents(title string, category string, city string, date string, tags []string) ([]*models.Event, error) {
	message := logMessage + "GetEvents:"
	log.Debug(message + "started")
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
		return nil, err
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			log.Error(message, "err = ", err)
			return nil, error2.ErrPostgres
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}
	log.Debug(message + "ended")
	return resultEvents, nil
}

func (s *Repository) GetVisitedEvents(userId string) ([]*models.Event, error) {
	message := logMessage + "GetVisitedEvents:"
	log.Debug(message + "started")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query := visitedQuery
	rows, err := s.db.Queryx(query, userIdInt)
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
	log.Debug(message + "ended")
	return resultEvents, nil
}

func (s *Repository) GetCreatedEvents(authorId string) ([]*models.Event, error) {
	message := logMessage + "GetCreatedEvents:"
	log.Debug(message + "started")
	authorIdInt, err := strconv.Atoi(authorId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query := createdQuery
	rows, err := s.db.Queryx(query, authorIdInt)
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
	log.Debug(message + "ended")
	return resultEvents, nil
}

func (s *Repository) Visit(eventId string, userId string) error {
	message := logMessage + "Visit:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := visitQuery
	_, err = s.db.Query(query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) Unvisit(eventId string, userId string) error {
	message := logMessage + "Unvisit:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := unvisitQuery
	_, err = s.db.Query(query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) IsVisited(eventId string, userId string) (bool, error) {
	message := logMessage + "IsVisited:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return false, error2.ErrAtoi
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return false, error2.ErrAtoi
	}
	query := isVisitedQuery
	var count int
	result := false
	err = s.db.Get(&count, query, eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return false, error2.ErrPostgres
	}
	if count > 0 {
		result = true
	}
	log.Debug(message + "ended")
	return result, nil
}

func (s *Repository) GetCities() ([]string, error) {
	message := logMessage + "GetCities:"
	log.Debug(message + "started")
	query := getCitiesQuery
	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Error(err)
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultCities []string
	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			log.Error(err)
			return nil, error2.ErrPostgres
		}
		resultCities = append(resultCities, c)
	}
	log.Debug(message + "ended")
	return resultCities, nil
}
