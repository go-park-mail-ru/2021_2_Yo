package postgres

import (
	"backend/auth"
	"backend/models"
	sql "github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
	log.Info("Here")
	newUser := toPostgresUser(user)
	insertQuery :=
		`insert into "user" (name, surname, mail, password) values($1, $2, $3, $4)`
	_, err := s.db.Exec(insertQuery, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password)
	if err != nil {
		log.Println("Auth:Repository:Postgres:CreateUser err :", err)
		return err
	}
	return nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	query := `select * from "user" where mail = $1 and password = $2`
	user := User{}
	err := s.db.QueryRow(query, mail, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password)
	if err != nil {
		log.Error("PostgresRepo : GetUser : err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	query := `select * from "user" where id = $1`
	user := User{}
	err := s.db.QueryRow(query, userId).Scan(&user.ID, &user.Name, &user.Surname, &user.Mail, &user.Password)
	if err != nil {
		log.Error("PostgresRepo : GetUser : err =", err)
		return nil, auth.ErrUserNotFound
	}
	return toModelUser(&user), nil
}
