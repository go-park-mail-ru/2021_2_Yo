package postgres

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	sql "github.com/jmoiron/sqlx"
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

func (s *Repository) CreateUser(user *models.User) error {
	message := logMessage + "CreateUser:"
	newUser := toPostgresUser(user)
	insertQuery := createUserQuery
	//TODO: Выяснить, нужен ли фронту user.id
	log.Debug(message+"password =", user.Password)
	_, err := s.db.Exec(insertQuery, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About)
	if err != nil {
		log.Debug(message+"err = ", err)
		return err
	}
	return nil
}

func (s *Repository) UpdateUser(user *models.User) error {
	message := "UpdateUser"
	newUser := toPostgresUser(user)
	insertQuery := `update "user" set name = $1, surname = $2, mail = $3, password = $4, about = $5 where id = $6`
	_, err := s.db.Exec(insertQuery, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About, newUser.ID)
	if err != nil {
		log.Debug(message+"err = ", err)
		return err
	}
	return nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	message := "GetUser"
	query := getUserQuery
	user := User{}
	err := s.db.QueryRow(query, mail, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password, &user.About)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	message := "GetUserById"
	query := getUserByIdQuery
	user := User{}
	err := s.db.QueryRow(query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password, &user.About)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}
