package repository

import (
	"backend/session/models"
	"errors"
	"fmt"
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
		return err
	}
	if result != "OK" {
		err := errors.New("Creating session ID error")
		return err
	}
	return nil
}

func (s *Repository) Check(sessionId string) (string, error) {
	result, err := redis.String(s.db.Do("GET", sessionId))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *Repository) Delete(sessionId string) error {
	result, err := redis.String(s.db.Do("DEL", sessionId))
	if err != nil {
		return err
	}
	fmt.Println("DEBUG!!! Session result = ", result)
	return nil
}
