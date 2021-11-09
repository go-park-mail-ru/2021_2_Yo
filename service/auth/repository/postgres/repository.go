package postgres

import (
	log "backend/logger"
	"backend/models"
	error2 "backend/service/auth/error"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage      = "service:auth:repository:postgres:"
	createUserQuery = `insert into "user" (name, surname, mail, password, about, img_url) values($1, $2, $3, $4, $5, $6) returning id`
	getUserQuery    = `select * from "user" where mail = $1 and password = $2`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

//TODO: ErrDataExists

func (s *Repository) CreateUser(user *models.User) (string, error) {
	newUser := toPostgresUser(user)
	log.Debug(logMessage+"CreateUser:newUser =", newUser)
	query := createUserQuery
	var userId int
	err := s.db.Get(&userId, query, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About, newUser.ImgUrl)
	if err != nil {
		return "", error2.ErrPostgres
	}
	return strconv.Itoa(userId), nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	query := getUserQuery
	user := User{}
	err := s.db.Get(&user, query, mail, password)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, error2.ErrUserNotFound
		}
		return nil, error2.ErrPostgres
	}
	return toModelUser(&user), nil
}
