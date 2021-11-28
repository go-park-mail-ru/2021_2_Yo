package postgres

import (
	log "backend/pkg/logger"
	"backend/pkg/models"
	error2 "backend/service/user/error"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage                       = "microservice:user:repository:"
	getUserByIdQuery                 = `select * from "user" where id = $1`
	updateUserInfoQueryWithoutImgUrl = `update "user" set name = $1, surname = $2, about = $3 where id = $4`
	updateUserInfoQuery              = `update "user" set name = $1, surname = $2, about = $3, img_url = $4 where id = $5`
	updateUserPasswordQuery          = `update "user" set password = $1 where id = $2`
	//TODO: updateUserImg в отдельный метод
	updateUserImgUrlQuery = `update "user" set img_url = $1 where id = $2`
	getSubscribersQuery   = `select u.* from "user" as u join subscribe s on s.subscriber_id = u.id where s.subscribed_id = $1`
	getSubscribesQuery    = `select u.* from "user" as u join subscribe s on s.subscribed_id = u.id where s.subscriber_id = $1`
	getVisitorsQuery      = `select u.* from "user" as u join visitor v on u.id = v.user_id where v.event_id = $1`
	subscribeQuery        = `insert into "subscribe" (subscribed_id, subscriber_id) values ($1, $2)`
	unsubscribeQuery      = `delete from subscribe where subscribed_id = $1 and subscriber_id = $2`
	isSubscribedQuery     = `select count(*) from subscribe where subscribed_id = $1 and subscriber_id = $2`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

/*
	GetUserById(userId string) (*models.User, error)
	///////
	UpdateUserInfo(user *models.User) error
	UpdateUserPassword(userId string, password string) error
	///////
	GetSubscribers(userId string) ([]*models.User, error)
	GetSubscribes(userId string) ([]*models.User, error)
	GetVisitors(eventId string) ([]*models.User, error)
	///////
	Subscribe(subscribedId string, subscriberId string) error
	Unsubscribe(subscribedId string, subscriberId string) error
	IsSubscribed(subscribedId string, subscriberId string) (bool, error)
*/

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	message := logMessage + "GetUserById:"
	log.Debug(message + "started")
	log.Debug(message + "userId = ", userId)
	query := getUserByIdQuery
	user := &User{}
	log.Debug(message+"HERE 1")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	err = s.db.Get(user, query, userIdInt)
	log.Debug(message+"HERE 2")
	if err != nil {
		log.Debug(message+"err = ",err)
		if err == sql2.ErrNoRows {
			return nil, error2.ErrUserNotFound
		}
		return nil, error2.ErrPostgres
	}
	log.Debug(message+"HERE 3")
	log.Debug(message + "user = ", user)
	log.Debug(message+"HERE 4")
	modelUser := toModelUser(user)
	log.Debug(message + "modelUser = ", modelUser)
	log.Debug(message + "ended")
	return modelUser, nil
}

func (s *Repository) UpdateUserInfo(u *models.User) error {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	postgresUser, err := toPostgresUser(u)
	if err != nil {
		return err
	}
	var query string
	if postgresUser.ImgUrl == "" {
		query = updateUserInfoQueryWithoutImgUrl
		_, err = s.db.Query(query, postgresUser.Name, postgresUser.Surname, postgresUser.About, postgresUser.ID)
		if err != nil {
			return error2.ErrPostgres
		}
	} else {
		query = updateUserInfoQuery
		_, err = s.db.Query(query, postgresUser.Name, postgresUser.Surname, postgresUser.About, postgresUser.ImgUrl, postgresUser.ID)
		if err != nil {
			return error2.ErrPostgres
		}
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) UpdateUserPassword(userId string, password string) error {
	message := logMessage + "UpdateUserPassword:"
	log.Debug(message + "started")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := updateUserPasswordQuery
	_, err = s.db.Query(query, password, userIdInt)
	if err != nil {
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) GetSubscribers(userId string) ([]*models.User, error) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query := getSubscribersQuery
	rows, err := s.db.Queryx(query, userIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultUsers []*models.User
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		modelUser := toModelUser(&u)
		resultUsers = append(resultUsers, modelUser)
	}
	log.Debug(message + "ended")
	return resultUsers, nil
}

func (s *Repository) GetSubscribes(userId string) ([]*models.User, error) {
	message := logMessage + "GetSubscribes:"
	log.Debug(message + "started")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query := getSubscribesQuery
	rows, err := s.db.Queryx(query, userIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultUsers []*models.User
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		modelUser := toModelUser(&u)
		resultUsers = append(resultUsers, modelUser)
	}
	log.Debug(message + "ended")
	return resultUsers, nil
}

func (s *Repository) GetVisitors(eventId string) ([]*models.User, error) {
	message := logMessage + "GetVisitors:"
	log.Debug(message + "started")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		return nil, error2.ErrAtoi
	}
	query := getVisitorsQuery
	rows, err := s.db.Queryx(query, eventIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}
	defer rows.Close()
	var resultUsers []*models.User
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, error2.ErrPostgres
		}
		modelUser := toModelUser(&u)
		resultUsers = append(resultUsers, modelUser)
	}
	log.Debug(message + "ended")
	return resultUsers, nil
}

func (s *Repository) Subscribe(subscribedId string, subscriberId string) error {
	message := logMessage + "Subscribe:"
	log.Debug(message + "started")
	subscribedIdInt, err := strconv.Atoi(subscribedId)
	if err != nil {
		return error2.ErrAtoi
	}
	subscriberIdInt, err := strconv.Atoi(subscriberId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := subscribeQuery
	_, err = s.db.Query(query, subscribedIdInt, subscriberIdInt)
	if err != nil {
		log.Error(message+"err = ", err)
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) Unsubscribe(subscribedId string, subscriberId string) error {
	message := logMessage + "Unsubscribe:"
	log.Debug(message + "started")
	subscribedIdInt, err := strconv.Atoi(subscribedId)
	if err != nil {
		return error2.ErrAtoi
	}
	subscriberIdInt, err := strconv.Atoi(subscriberId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := unsubscribeQuery
	_, err = s.db.Query(query, subscribedIdInt, subscriberIdInt)
	if err != nil {
		log.Error(err)
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) IsSubscribed(subscribedId string, subscriberId string) (bool, error) {
	message := logMessage + "IsSubscribed:"
	log.Debug(message + "started")
	subscribedIdInt, err := strconv.Atoi(subscribedId)
	if err != nil {
		return false, error2.ErrAtoi
	}
	subscriberIdInt, err := strconv.Atoi(subscriberId)
	if err != nil {
		return false, error2.ErrAtoi
	}
	query := isSubscribedQuery
	var count int
	result := false
	err = s.db.Get(&count, query, subscribedIdInt, subscriberIdInt)
	if err != nil {
		return false, error2.ErrPostgres
	}
	if count > 0 {
		result = true
	}
	log.Debug(message + "ended")
	return result, nil
}
