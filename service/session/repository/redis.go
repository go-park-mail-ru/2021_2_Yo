package repository

import (
	log "backend/logger"
	"backend/service/session/models"
	"github.com/go-redis/redis"
)

type Repository struct {
	db *redis.Client
}

func NewRepository(database *redis.Client) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) Create(data *models.SessionData) error {
	res := s.db.Set(data.SessionId, data.UserId, data.Expiration)
	return res.Err()
}

func (s *Repository) Check(sessionId string) (string, error) {
	res := s.db.Get(sessionId)
	return res.Val(), res.Err()
}

func (s *Repository) Delete(sessionId string) error {
	res := s.db.Del(sessionId)
	log.Error("Redis delete err =", res.Err())
	log.Debug("Redis delete result = ", res.Val())
	return res.Err()
}
