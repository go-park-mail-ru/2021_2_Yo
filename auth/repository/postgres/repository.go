package postgres

import (
	"backend/auth"
	"backend/models"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage       = "auth:repository:postgres:"
	createUserQuery  = `insert into "user" (name, surname, mail, password, about) values($1, $2, $3, $4, $5) returning id`
	updateUserQuery  = ``
	getUserQuery     = `select * from "user" where mail = $1 and password = $2`
	getUserByIdQuery = `select * from "user" where id = $1`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) CreateUser(user *models.User) (string, error) {
	newUser := toPostgresUser(user)
	query := createUserQuery
	var userId int
	err := s.db.QueryRow(query, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About).Scan(&userId)
	if err != nil {
		return "", auth.ErrPostgres
	}
	return strconv.Itoa(userId), nil
}

func (s *Repository) UpdateUserInfo(userId, name, surname, about string) error {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return auth.ErrAtoi
	}
	query := `update "user" set name = $1, surname = $2, about = $3 where id = $4`
	_, err = s.db.Exec(query, name, surname, about, userIdInt)
	if err != nil {
		return auth.ErrPostgres
	}
	return nil
}

func (s *Repository) UpdateUserPassword(userId, password string) error {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return auth.ErrAtoi
	}
	query := `update "user" set password = $1 where id = $2`
	_, err = s.db.Exec(query, password, userIdInt)
	if err != nil {
		return auth.ErrPostgres
	}
	return nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	query := getUserQuery
	user := User{}
	err := s.db.QueryRow(query, mail, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password, &user.About)
	if err == sql2.ErrNoRows {
		return nil, auth.ErrUserNotFound
	}
	if err != nil {
		return nil, auth.ErrPostgres
	}
	return toModelUser(&user), nil
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	query := getUserByIdQuery
	user := User{}
	err := s.db.QueryRow(query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password, &user.About)
	if err == sql2.ErrNoRows {
		return nil, auth.ErrUserNotFound
	}
	if err != nil {
		return nil, auth.ErrPostgres
	}
	return toModelUser(&user), nil
}
