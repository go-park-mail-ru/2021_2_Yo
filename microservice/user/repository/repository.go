package repository

import (
	proto "backend/microservice/user/proto"
	log "backend/pkg/logger"
	"backend/pkg/models"
	error2 "backend/service/user/error"
	"context"
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
	subscribeQuery        = `insert into "subscribe" (subscribed_id, subscriber_id) values ($1, $2)`
	getSubscribersQuery   = `select u.* from "user" as u join subscribe s on s.subscriber_id = u.id where s.subscribed_id = $1`
	getSubscribesQuery    = `select u.* from "user" as u join subscribe s on s.subscribed_id = u.id where s.subscriber_id = $1`
	getVisitorsQuery      = `select u.* from "user" as u join visitor v on u.id = v.user_id where v.event_id = $1`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) GetUserById(ctx context.Context, in *proto.UserId) (*proto.User, error) {
	message := logMessage + "GetUserById:"
	log.Debug(message + "started")

	userId := in.ID

	query := getUserByIdQuery
	user := User{}
	err := s.db.Get(&user, query, userId)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrUserNotFound
		}
		return nil, error2.ErrPostgres
	}

	modelUser := toModelUser(&user)
	protoUser := toProtoUser(modelUser)

	log.Debug(message + "ended")
	return protoUser, nil

}

func (s *Repository) UpdateUserInfo(ctx context.Context, in *proto.User) (*proto.Empty, error) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")

	postgresUser, err := toPostgresUser(&models.User{
		ID:       in.ID,
		Name:     in.Name,
		Surname:  in.Surname,
		Mail:     in.Mail,
		Password: in.Password,
		About:    in.About,
		ImgUrl:   in.ImgUrl,
	})
	if err != nil {
		return nil, err
	}

	var query string
	if postgresUser.ImgUrl == "" {
		query = updateUserInfoQueryWithoutImgUrl
		_, err = s.db.Query(query, postgresUser.Name, postgresUser.Surname, postgresUser.About, postgresUser.ID)
		if err != nil {
			return nil, error2.ErrPostgres
		}
	} else {
		query = updateUserInfoQuery
		_, err = s.db.Query(query, postgresUser.Name, postgresUser.Surname, postgresUser.About, postgresUser.ImgUrl, postgresUser.ID)
		if err != nil {
			return nil, error2.ErrPostgres
		}
	}

	log.Debug(message + "ended")
	return nil, nil

}

func (s *Repository) UpdateUserPassword(ctx context.Context, in *proto.UpdateUserPasswordRequest) (*proto.Empty, error) {

	message := logMessage + "UpdateUserPassword:"
	log.Debug(message + "started")

	userId := in.ID
	password := in.Password

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	query := updateUserPasswordQuery
	_, err = s.db.Query(query, password, userIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}

	log.Debug(message + "ended")
	return nil, nil

}

func (s *Repository) Subscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.Empty, error) {

	message := logMessage + "Subscribe:"
	log.Debug(message + "started")

	subscribedId := in.SubscribedId
	subscriberId := in.SubscriberId

	subscribedIdInt, err := strconv.Atoi(subscribedId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	subscriberIdInt, err := strconv.Atoi(subscriberId)
	if err != nil {
		return nil, error2.ErrAtoi
	}

	query := subscribeQuery
	_, err = s.db.Query(query, subscribedIdInt, subscriberIdInt)
	if err != nil {
		return nil, error2.ErrPostgres
	}

	log.Debug(message + "ended")
	return nil, nil

}

func (s *Repository) GetSubscribers(ctx context.Context, in *proto.UserId) (*proto.Users, error) {

	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")

	userId := in.ID

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

	outUsers := make([]*proto.User, len(resultUsers))
	for i, event := range resultUsers {
		outUsers[i] = toProtoUser(event)
	}
	out := &proto.Users{Users: outUsers}

	log.Debug(message + "ended")
	return out, nil
}

func (s *Repository) GetSubscribes(ctx context.Context, in *proto.UserId) (*proto.Users, error) {

	message := logMessage + "GetSubscribes:"
	log.Debug(message + "started")

	userId := in.ID

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

	outUsers := make([]*proto.User, len(resultUsers))
	for i, event := range resultUsers {
		outUsers[i] = toProtoUser(event)
	}
	out := &proto.Users{Users: outUsers}

	log.Debug(message + "ended")
	return out, nil

}

func (s *Repository) GetVisitors(ctx context.Context, in *proto.EventId) (*proto.Users, error) {

	message := logMessage + "GetVisitors:"
	log.Debug(message + "started")

	eventId := in.ID

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

	outUsers := make([]*proto.User, len(resultUsers))
	for i, event := range resultUsers {
		outUsers[i] = toProtoUser(event)
	}
	out := &proto.Users{Users: outUsers}

	log.Debug(message + "ended")
	return out, nil

}
