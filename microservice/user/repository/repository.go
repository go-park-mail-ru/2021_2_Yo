package repository

import (
	proto "backend/microservice/user/proto"
	"backend/models"
	error2 "backend/service/user/error"
	"context"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage                       = "service:user:repository:postgres:"
	getUserByIdQuery                 = `select * from "user" where id = $1`
	updateUserInfoQueryWithoutImgUrl = `update "user" set name = $1, surname = $2, about = $3 where id = $4`
	updateUserInfoQuery              = `update "user" set name = $1, surname = $2, about = $3, img_url = $4 where id = $5`
	updateUserPasswordQuery          = `update "user" set password = $1 where id = $2`
	//TODO: updateUserImg в отдельный метод
	updateUserImgUrl = `update "user" set img_url = $1 where id = $2`
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

	return protoUser, nil

}

func (s *Repository) UpdateUserInfo(ctx context.Context, in *proto.User) (*proto.Empty, error) {

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

	return nil, nil

}

func (s *Repository) UpdateUserPassword(ctx context.Context, in *proto.UpdateUserPasswordRequest) (*proto.Empty, error) {

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

	return nil, nil

}