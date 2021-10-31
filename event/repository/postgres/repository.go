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
	checkAuthorQuery = `select author_id from "event" where id = $1`
	listQuery        = `select * from "event"`
	getEventQuery    = `select * from "event" where id = $1`
	createEventQuery = `insert into "event" 
		(title, description, text, city, category, viewed, img_url, date, geo, author_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		returning id`
	updateEventQuery = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5, viewed = $6, img_url = $7, date = $8, geo = $9 
		where event.id = $10`
	deleteEventQuery = `delete from "event" where event.author_id = $1`
)

func (s *Repository) checkAuthor(eventId int, userId int) (bool, error) {
	var authorId int
	query := checkAuthorQuery
	err := s.db.QueryRow(query, userId).Scan(&authorId)
	if err != nil {
		log.Debug("checkAuthor err1 = ", err)
		return false, err
	}
	log.Debug("checkAuthor HERE")
	if authorId == userId {
		return true, nil
	} else {
		return false, nil
	}
}

/*
func (s *Repository) getTagsIds(eventId int) ([]int, error) {
	message := logMessage + "getTagsIds:"
	query :=
}
*/

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
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return nil, event.ErrEventNotFound
}

func (s *Repository) CreateEvent(e *models.Event) (string, error) {
	message := logMessage + "CreateEvent:"
	newEvent, err := toPostgresEvent(e)
	if err != nil {
		return "", err
	}
	var eventId int
	query := createEventQuery
	err = s.db.QueryRow(query,
		newEvent.Title,
		newEvent.Description,
		newEvent.Text,
		newEvent.City,
		newEvent.Category,
		newEvent.Viewed,
		newEvent.Img_Url,
		newEvent.Date,
		newEvent.Geo,
		newEvent.Author_ID).Scan(&eventId)
	if err != nil {
		log.Debug(message+"err = ", err)
		return "", err
	}
	return strconv.Itoa(eventId), nil
}

func (s *Repository) UpdateEvent(updatedEvent *models.Event, userId string) error {
	message := logMessage + "UpdateEvent:"
	e, err := toPostgresEvent(updatedEvent)
	if err != nil {
		return err
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	canUpdate, err := s.checkAuthor(e.ID, userIdInt)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	if !canUpdate {
		err = errors.New("user can't update event")
		log.Error(message+"err =", err)
		return err
	}
	log.Debug(message + "HERE")
	query := updateEventQuery
	_, err = s.db.Exec(query, e.Title, e.Description, e.Text, e.City, e.Category, e.Viewed, e.Img_Url, e.Date, e.Geo, e.ID)
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
	query := deleteEventQuery
	_, err = s.db.Exec(query, userIdInt)
	if err != nil {
		log.Debug(message+"err = ", err)
		return err
	}
	return nil
}
