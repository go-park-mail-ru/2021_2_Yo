package postgres

import (
	models2 "backend/internal/models"
	error2 "backend/internal/service/event/error"
	log "backend/pkg/logger"
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
	logMessage          = "service:event:repository:postgres:"
	checkAuthorQuery    = `select author_id from "event" where id = $1`
	listQuery           = `select * from "event"`
	incrementEventViews = `update "event" set viewed = viewed + 1 where event.id = $1`
	getEventQuery       = `select * from "event" where id = $1`
	createEventQuery    = `insert into "event" 
		(title, description, text, city, category, viewed, img_url, date, geo, address, tag, author_id) 
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::varchar[], $12) 
		returning id`
	updateEventQuery = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5,
		img_url = $6, date = $7, geo = $8, address = $9, tag = $10
		where event.id = $11`
	updateEventQueryWithoutImgUrl = `update "event" set
		title = $1, description = $2, text = $3, city = $4, category = $5,
		date = $6, geo = $7, address = $8, tag = $9 
		where event.id = $10`
	deleteEventQuery = `delete from "event" where id = $1`
	visitedQuery     = `select e.* from "event" as e join visitor as v on v.event_id = e.id where v.user_id = $1`
	createdQuery     = `select * from "event" where author_id = $1`
	visitQuery       = `insert into "visitor" (event_id, user_id) values ($1, $2)`
	unvisitQuery     = `delete from "visitor" where event_id = $1 and user_id = $2`
	isVisitedQuery   = `select count(*) from "visitor" where event_id = $1 and user_id = $2`
	getCitiesQuery   = `select distinct city from event`
	getSubsInfo      = `select u2.name, u2.mail, e.title, e.img_url from "user" as u1 join subscribe on u1.id = subscribe.subscribed_id
							join "user" as u2 on u2.id = subscribe.subscriber_id 
							join "event" as e on e.id = $1
							where e.author_id = u1.id`
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

func (s *Repository) CreateEvent(e *models2.Event) (string, error) {
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

func (s *Repository) UpdateEvent(e *models2.Event, userId string) error {
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
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) GetEventById(eventId string) (*models2.Event, error) {
	message := logMessage + "GetEventById:"
	log.Debug(message + "started")
	var query string
	var e Event
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query = incrementEventViews
	s.db.Query(query, eventIdInt)
	query = getEventQuery
	err = s.db.Get(&e, query, eventIdInt)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrNoRows
		}
		return nil, error2.ErrPostgres
	}
	modelEvent := toModelEvent(&e)
	log.Debug(message + "ended")
	return modelEvent, nil
}

func (s *Repository) GetEvents(userId string, title string, category string, city string, date string, tags []string) ([]*models2.Event, error) {
	message := logMessage + "GetEvents:"
	log.Debug(message + "started")
	postgresTags := make(pq.StringArray, len(tags))

	var userIdInt int

	for i := range tags {
		postgresTags[i] = tags[i]
	}
	query := `select e.*, count(v) from event as e
				left join visitor as v on e.id = v.event_id and `
	query += `v.user_id = $1 `
	if userId == "" {
		userIdInt = 0
	} else {
		userIdInt1, err := strconv.Atoi(userId)
		if err != nil {
			return nil, error2.ErrAtoi
		}
		userIdInt = userIdInt1
	}
	if title != "" {
		query += `where lower(title) ~ lower($2) and `
	} else {
		query += `where $2 = $2 and `
	}
	if category != "" {
		query += `lower(category) = lower($3) and `
	} else {
		query += `$3 = $3 and `
	}
	if city != "" {
		query += `lower(city) = lower($4) and `
	} else {
		query += `$4 = $4 and `
	}
	if date != "" {
		query += `lower(e.date) = lower($5) and `
	} else {
		query += `$5 = $5 and `
	}
	if len(postgresTags) != 0 {
		query += `tag && $6::varchar[] `
	} else {
		query += `$6 = $6 `
	}
	query += `group by e.id,
         e.title,
         e.description,
         e.text,
         e.city,
         e.category,
         e.viewed,
         e.img_url,
         e.date,
         e.geo,
         e.address,
         e.tag,
         e.author_id 
         order by viewed DESC`
	rows, err := s.db.Queryx(query, userIdInt, title, category, city, date, postgresTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resultEvents []*models2.Event
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

func (s *Repository) GetVisitedEvents(userId string) ([]*models2.Event, error) {
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
	var resultEvents []*models2.Event
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

func (s *Repository) GetCreatedEvents(authorId string) ([]*models2.Event, error) {
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
	var resultEvents []*models2.Event
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
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultCities []string
	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		resultCities = append(resultCities, c)
	}
	log.Debug(message + "ended")
	return resultCities, nil
}

func (s *Repository) EmailNotify(eventId string) ([]*models2.Info, error) {
	message := logMessage + "EmailNotify:"
	log.Debug(message + "started")
	query := getSubsInfo
	rows, err := s.db.Queryx(query, eventId)
	if err != nil {
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var subsInfo []*models2.Info
	for rows.Next() {
		var mail, name, title, img_url string
		err := rows.Scan(&name, &mail, &title, &img_url)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		userInfo := &models2.Info{
			Name:    name,
			Mail:    mail,
			Title:   title,
			Img_url: img_url,
		}
		subsInfo = append(subsInfo, userInfo)
	}
	log.Debug(message + "ended")
	return subsInfo, nil
}
