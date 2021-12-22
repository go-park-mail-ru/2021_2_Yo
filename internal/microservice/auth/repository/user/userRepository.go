package user

import (
	"backend/internal/models"
	error2 "backend/internal/service/auth/error"
	log "backend/pkg/logger"
	sql2 "database/sql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

const (
	logMessage      = "service:auth:repository:postgres:"
	createUserQuery = `insert into "user" (name, surname, mail, password, about) values($1, $2, $3, $4, $5) returning id`
	getUserQuery    = `select * from "user" where mail = $1 and password = $2`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) CreateUser(user *models.User) (string, error) {
	newUser := toPostgresUser(user)
	query := createUserQuery
	var userId int
	err := s.db.Get(&userId, query, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			return "", error2.ErrUserExists
		}
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
