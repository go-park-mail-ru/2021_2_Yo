package repository

import (
	error2 "backend/session/error"
	"backend/session/models"
	"github.com/gomodule/redigo/redis"
)

type Repository struct {
	db redis.Conn
}

func NewRepository(database redis.Conn) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) Create(data *models.SessionData) error {
	result, err := redis.String(s.db.Do("SET", data.SessionId, data.UserId, "EX", data.Expiration))
	if err != nil {
		return error2.ErrRedis
	}
	if result != "OK" {
		return error2.ErrCreateSession
	}
	return nil
}

func (s *Repository) Check(sessionId string) (string, error) {
	result, err := redis.String(s.db.Do("GET", sessionId))
	if err != nil {
		return "", error2.ErrCheckSession
	}
	return result, nil
}

func (s *Repository) Delete(sessionId string) error {
	_, err := redis.String(s.db.Do("DEL", sessionId))
	if err != nil {
		return error2.ErrDeleteSession
	}
	return nil
}
