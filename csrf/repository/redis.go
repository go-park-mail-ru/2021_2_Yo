package repository

import (
	"backend/csrf/models"
	"github.com/gomodule/redigo/redis"
	my_error "backend/csrf/error"
)

type Repository struct {
	db redis.Conn
}

func NewRepository(database redis.Conn) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) Create(data *models.CSRFData) error {
	result, err := redis.String(s.db.Do("SET", data.CSRFToken, data.UserId, "EX", data.Expiration))
	if err != nil {
		return my_error.ErrRedis
	}
	if result != "OK" {
		return my_error.ErrRedis
	}
	return nil
}

func (s *Repository) Check(susToken string) (string, error) {
	result, err := redis.String(s.db.Do("GET", susToken))
	if err != nil {
		return "", my_error.ErrRedis
	}
	return result, nil
}

func (s *Repository) Delete(susToken string) error {
	_, err := redis.String(s.db.Do("DEL", susToken))
	if err != nil {
		return my_error.ErrRedis
	}
	return nil
}