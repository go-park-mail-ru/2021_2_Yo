package postgres

import (
	"backend/auth"
	"backend/models"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	result := &Repository{
		db: database,
	}
	return result
}

func (s *Repository) CreateUser(user *models.User) error {
	newUser := toPostgresUser(user)
	insertQuery :=
		`insert into users (name, surname, mail, password) values($1, $2, $3, $4)`
	_, err := s.db.Exec(insertQuery, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password)
	if err != nil {
		log.Println("Auth:Repository:Postgres:CreateUser err :", err)
		return err
	}
	return nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	query := `select * from users where mail = $1 and password = $2`
	user := User{}
	err := s.db.QueryRow(query, mail, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password)
	if err != nil {
		log.Error("PostgresRepo : GetUser : err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	query := `select * from users where id = $1`
	user := User{}
	err := s.db.QueryRow(query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password)
	if err != nil {
		log.Error("PostgresRepo : GetUser : err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}
