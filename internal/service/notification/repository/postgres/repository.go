package postgres

import (
	"backend/internal/models"
	error2 "backend/internal/service/notification/error"
	log "backend/pkg/logger"
	sql "github.com/jmoiron/sqlx"
)

const (
	logMessage = "service:notification:repository:postgres:"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) CreateSubscribeNotification(receiverId string, user *models.User, seen bool) error {
	message := logMessage + "CreateSubNotification:"
	log.Debug(message + "started")
	query := `insert into "notification" (type, receiver_id, user_id, user_name, user_surname, user_img_url, seen) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Query(query, "sub", receiverId, user.ID, user.Name, user.Surname, user.ImgUrl, seen)
	if err != nil {
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) CreateInviteNotification(receiverId string, user *models.User, event *models.Event, seen bool) error {
	message := logMessage + "CreateInvNotification:"
	log.Debug(message + "started")
	query := `insert into "notification" (type, receiver_id, user_id, user_name, user_surname, user_img_url, event_id, event_title, seen) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.Query(query, "inv", receiverId, user.ID, user.Name, user.Surname, user.ImgUrl, event.ID, event.Title, seen)
	if err != nil {
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) CreateNewEventNotification(receiverId string, user *models.User, event *models.Event, seen bool) error {
	message := logMessage + "CreateInvNotification:"
	log.Debug(message + "started")
	query := `insert into "notification" (type, receiver_id, user_id, user_name, user_surname, user_img_url, event_id, event_title, seen) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.Query(query, "new", receiverId, user.ID, user.Name, user.Surname, user.ImgUrl, event.ID, event.Title, seen)
	if err != nil {
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) UpdateNotificationsStatus(userId string) error {
	message := logMessage + "UpdateNotificationsStatus:"
	log.Debug(message + "started")
	query := `update "notification" set seen = true where receiver_id = $1`
	_, err := s.db.Query(query, userId)
	if err != nil {
		return error2.ErrPostgres
	}
	log.Debug(message + "ended")
	return nil
}

func (s *Repository) GetAllNotifications(userId string) ([]*models.Notification, error) {
	message := logMessage + "GetAllNotifications:"
	log.Debug(message + "started")
	query := `select * from notification where receiver_id = $1`
	rows, err := s.db.Queryx(query, userId)
	if err != nil {
		log.Error(message, "err = ", err)
		return nil, err
	}
	defer rows.Close()
	var resultNotifications []*models.Notification
	for rows.Next() {
		var n Notification
		err := rows.StructScan(&n)
		if err != nil {
			log.Error(message, "err = ", err)
			return nil, error2.ErrPostgres
		}
		modelNotification := toModelNotification(&n)
		resultNotifications = append(resultNotifications, modelNotification)
	}
	log.Debug(message + "ended")
	return resultNotifications, nil
}

func (s *Repository) GetNewNotifications(userId string) ([]*models.Notification, error) {
	message := logMessage + "GetAllNotifications:"
	log.Debug(message + "started")
	query := `select * from notification where receiver_id = $1 and seen = false`
	rows, err := s.db.Queryx(query, userId)
	if err != nil {
		log.Error(message, "err = ", err)
		return nil, err
	}
	defer rows.Close()
	var resultNotifications []*models.Notification
	for rows.Next() {
		var n Notification
		err := rows.StructScan(&n)
		if err != nil {
			log.Error(message, "err = ", err)
			return nil, error2.ErrPostgres
		}
		modelNotification := toModelNotification(&n)
		resultNotifications = append(resultNotifications, modelNotification)
	}
	log.Debug(message + "ended")
	return resultNotifications, nil
}
