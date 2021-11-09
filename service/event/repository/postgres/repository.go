package postgres

import (
	log "backend/logger"
	"backend/models"
	error2 "backend/service/event/error"
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
	newEvent, err := toPostgresEvent(e)
	if err != nil {
		return "", err
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
			return "", error2.ErrNoRows
		}
		return "", error2.ErrPostgres
	}
	return strconv.Itoa(eventId), nil
}

func (s *Repository) UpdateEvent(updatedEvent *models.Event, userId string) error {
	eventIdInt, err := strconv.Atoi(updatedEvent.ID)
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
	e, err := toPostgresEvent(updatedEvent)
	if err != nil {
		return err
	}
	e.ID = eventIdInt
	var query string
	if e.ImgUrl != "" {
		query = updateEventQuery
		_, err = s.db.Query(query,
			e.Title,
			e.Description,
			e.Text,
			e.City,
			e.Category,
			e.Viewed,
			e.ImgUrl,
			e.Date,
			e.Geo,
			e.Tag,
			e.ID)
		if err != nil {
			return error2.ErrPostgres
		}
	} else {
		query = updateEventQueryWithoutImgUrl
		_, err = s.db.Query(query,
			e.Title,
			e.Description,
			e.Text,
			e.City,
			e.Category,
			e.Viewed,
			e.Date,
			e.Geo,
			e.Tag,
			e.ID)
		if err != nil {
			return error2.ErrPostgres
		}
	}
	return nil
}

func (s *Repository) DeleteEvent(eventId string, userId string) error {
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
	_, err = s.db.Exec(query, eventIdInt)
	if err != nil {
		return error2.ErrPostgres
	}
	return nil
}

func (s *Repository) GetEventById(eventId string) (*models.Event, error) {
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
	var resultEvent *models.Event
	resultEvent = toModelEvent(&e)
	return resultEvent, nil
}

func (s *Repository) GetEvents(title string, category string, tags []string) ([]*models.Event, error) {
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
	return resultEvents, nil
}

func (s *Repository) GetEventsFromAuthor(authorId string) ([]*models.Event, error) {
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
	return resultEvents, nil
}
