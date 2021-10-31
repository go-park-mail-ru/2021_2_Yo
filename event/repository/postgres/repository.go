package postgres

import (
	"backend/event"
	log "backend/logger"
	"backend/models"
	"errors"
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

func (s *Repository) checkAuthor(eventId int, userId int) (bool, error) {
	query := `select author_id from "event" where id = $1`
	rows, err := s.db.Queryx(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		var authorId int
		err := rows.StructScan(&authorId)
		if err != nil {
			return false, err
		}
		if authorId == userId {
			return true, nil
		} else {
			return false, nil
		}
	}
	err = errors.New("Unexpected error")
	return false, err
}

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
	message := logMessage + "CreateEvent:"
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

func (s *Repository) UpdateEvent(updatedEvent *models.Event) error {
	message := logMessage + "UpdateEvent:"
	e, _ := toPostgresEvent(updatedEvent)
	query :=
		`update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, viewed = $6, img_url = $7, date = $8, geo = $9 
		where event.id = $10`
	_, err := s.db.Exec(query, e.Title, e.Description, e.Text, e.City, e.Category, e.Viewed, e.Img_Url, e.Date, e.Geo, e.ID)
	if err != nil {
		log.Debug(message+"err = ", err)
		return err
	}
	return nil
}

func (s *Repository) DeleteEvent(eventId string, userId string) error {
	message := logMessage + "DeleteEvent:"
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	canDelete, err := s.checkAuthor(eventIdInt, userIdInt)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	if !canDelete {
		err = errors.New("user can't delete event")
		log.Error(message+"err =", err)
		return err
	}
	query := `delete from "event" where event.author_id = $1`
	_, err = s.db.Exec(query, userIdInt)
	if err != nil {
		log.Debug(message+"err = ", err)
		return err
	}
	return nil
}
