package postgres

import (
	"backend/event"
	log "backend/logger"
	"backend/models"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const logMessage = "event:repository:postgres:"

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

const (
	listQuery     = `select * from "event"`
	getEventQuery = `select * from "event" where id = $1`
)

//TODO: Увеличить количество кастомных ошибок
//TODO: sql.ErrNoRows

func (s *Repository) List() ([]*models.Event, error) {
	message := logMessage + "List:"
	query := listQuery
	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			log.Error(message+"err =", err)
			return nil, err
		}
		modelEvent := toModelEvent(&e)
		resultEvents = append(resultEvents, modelEvent)
	}
	log.Debug(message+"resultEvents =", resultEvents)
	return resultEvents, nil
}

func (s *Repository) GetEvent(eventId string) (*models.Event, error) {
	message := logMessage + "GetEvent:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	query := getEventQuery
	rows, err := s.db.Queryx(query, eventIdInt)
	defer rows.Close()
	var resultEvent *models.Event
	if rows.Next() {
		var e Event
		err := rows.StructScan(&e)
		if err != nil {
			log.Error(message+"err =", err)
			return nil, err
		}
		resultEvent = toModelEvent(&e)
		return resultEvent, nil
	}
	return nil, event.ErrEventNotFound
}

func (s *Repository) CreateEvent(e *models.Event) (string, error) {
	message := "CreateEvent"
	newEvent, _ := toPostgresEvent(e)
	var eventId string
	query :=
		`insert into "event" 
		(title, description, text, city, category, viewed, img_url, date, geo, author_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		returning id`
	err := s.db.QueryRow(query,
		newEvent.Title,
		newEvent.Description,
		newEvent.Text,
		newEvent.City,
		newEvent.Category,
		newEvent.Viewed,
		newEvent.Img_Url,
		newEvent.Geo,
		newEvent.Author_ID).Scan(&eventId)
	if err != nil {
		log.Debug(message+"err = ", err)
		return "", err
	}
	return eventId, nil
}

func (s *Repository) UpdateEvent(eventId string, e *models.Event) error {
	return nil
}

//TODO: userId должен проверяться в useCase

func (s *Repository) DeleteEvent(eventId string, userId string) error {
	return nil
}
