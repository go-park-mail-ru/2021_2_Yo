package postgres

import (
	"backend/models"
	error2 "backend/service/user/error"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage                       = "service:user:repository:postgres:"
	getUserByIdQuery                 = `select * from "user" where id = $1`
	updateUserInfoQueryWithoutImgUrl = `update "user" set name = $1, surname = $2, about = $3, $4 = $4 where id = $5`
	updateUserInfoQuery              = `update "user" set name = $1, surname = $2, about = $3, img_url = $4 where id = $5`
	updateUserPasswordQuery          = `update "user" set password = $1 where id = $2`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	query := getUserByIdQuery
	user := User{}
	err := s.db.Get(&user, query, userId)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrUserNotFound
		}
		return nil, error2.ErrPostgres
	}
	return toModelUser(&user), nil
}

func (s *Repository) UpdateUserInfo(user *models.User) error {
	postgresUser, err := toPostgresUser(user)
	if err != nil {
		return err
	}
	var query string
	if postgresUser.ImgUrl == "" {
		query = updateUserInfoQueryWithoutImgUrl
	} else {
		query = updateUserInfoQuery
	}
	_, err = s.db.Query(query, postgresUser.Name, postgresUser.Surname, postgresUser.About, postgresUser.ImgUrl, postgresUser.ID)
	if err != nil {
		return error2.ErrPostgres
	}
	return nil
}

func (s *Repository) UpdateUserPassword(userId, password string) error {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := updateUserPasswordQuery
	_, err = s.db.Query(query, password, userIdInt)
	if err != nil {
		return error2.ErrPostgres
	}
	return nil
}
