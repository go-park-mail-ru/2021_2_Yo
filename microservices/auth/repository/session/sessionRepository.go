package session

import (
	log "backend/logger"
	"backend/service/session/models"
	"github.com/go-redis/redis"
)

const logMessage = "service:session:repository:"

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
	log.Debug(logMessage+"Create:res =", res)
	return res.Err()
}

func (s *Repository) Check(sessionId string) (string, error) {
	res := s.db.Get(sessionId)
	log.Debug("Check:res =", res)
	return res.Val(), res.Err()
}

func (s *Repository) Delete(sessionId string) error {
	res := s.db.Del(sessionId)
	log.Debug("Check:delete =", res)
	return res.Err()
}
