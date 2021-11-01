package postgres

import (
	error2 "backend/auth/error"
	"backend/models"
	sql2 "database/sql"
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	logMessage              = "auth:repository:postgres:"
	createUserQuery         = `insert into "user" (name, surname, mail, password, about) values($1, $2, $3, $4, $5) returning id`
	updateUserInfoQuery     = `update "user" set name = $1, surname = $2, about = $3 where id = $4`
	updateUserPasswordQuery = `update "user" set password = $1 where id = $2`
	getUserQuery            = `select * from "user" where mail = $1 and password = $2`
	getUserByIdQuery        = `select * from "user" where id = $1`
)

/*
type Event struct {
	ID          int
	Title       string
	Description string
	Text        string
	City        string
	Category    string
	Viewed      int
	ImgUrl      string `db:"img_url""`
	Tag         pq.StringArray
	Date        string
	Geo         string
	Author_ID   int
}

func main() {
	db, _ := InitPostgresDB()
	getEventQuery := `select * from "event" where id = $1`
	query := getEventQuery
	var e Event
	err := db.Get(&e, query, 1)
	if err != nil {
		log.Error(err)
	}
	log.Info(e.Tag)
	return
}
*/

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
	err := s.db.Get(&userId, query, newUser.Name, newUser.Surname, newUser.Mail, newUser.Password, newUser.About)
	if err != nil {
		return "", error2.ErrPostgres
	}
	return strconv.Itoa(userId), nil
}

func (s *Repository) UpdateUserInfo(userId, name, surname, about string) error {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return error2.ErrAtoi
	}
	query := updateUserInfoQuery
	_, err = s.db.Exec(query, name, surname, about, userIdInt)
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
	_, err = s.db.Exec(query, password, userIdInt)
	if err != nil {
		return error2.ErrPostgres
	}
	return nil
}

func (s *Repository) GetUser(mail, password string) (*models.User, error) {
	query := getUserQuery
	user := User{}
	err := s.db.Get(&user, query, mail, password)
	if err == sql2.ErrNoRows {
		return nil, error2.ErrUserNotFound
	}
	if err != nil {
		return nil, error2.ErrPostgres
	}
	return toModelUser(&user), nil
}

func (s *Repository) GetUserById(userId string) (*models.User, error) {
	query := getUserByIdQuery
	user := User{}
	err := s.db.Get(&user, query, userId)
	if err == sql2.ErrNoRows {
		return nil, error2.ErrUserNotFound
	}
	if err != nil {
		return nil, error2.ErrPostgres
	}
	return toModelUser(&user), nil
}
