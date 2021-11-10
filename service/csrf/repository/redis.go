package repository

import (
	"backend/service/csrf/models"
	"github.com/go-redis/redis"
	log "backend/logger"
)

const logMessage = "service:csrf:repository:"

type Repository struct {
	db *redis.Client
}

func NewRepository(database *redis.Client) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) Create(data *models.CSRFData) error {
	res := s.db.Set(data.CSRFToken, data.UserId, data.Expiration)
	log.Debug(logMessage+"Create:res =", res)
	return res.Err()
}

func (s *Repository) Check(susToken string) (string, error) {
	res := s.db.Get(susToken)
	log.Debug("Check:res =", res)
	return res.Val(), res.Err()
}

func (s *Repository) Delete(susToken string) error {
	res := s.db.Del(susToken)
	log.Debug("Check:delete =", res)
	return res.Err()
}
