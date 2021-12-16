package session

import (
	authServiceModels "backend/internal/microservice/auth/models"
	log "backend/pkg/logger"
	"github.com/go-redis/redis"
)

const logMessage = "service:session:repository:"

type Repository struct {
	db redis.Cmdable
}

func NewRepository(database redis.Cmdable) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) Create(data *authServiceModels.SessionData) error {
	message := logMessage + "Create:"
	log.Info(message + "started")
	res := s.db.Set(data.SessionId, data.UserId, data.Expiration)
	log.Info(message + "ended")
	return res.Err()
}

func (s *Repository) Check(sessionId string) (string, error) {
	message := logMessage + "Check:"
	log.Info(message + "started")
	res := s.db.Get(sessionId)
	log.Info(message + "ended")
	return res.Val(), res.Err()
}

func (s *Repository) Delete(sessionId string) error {
	message := logMessage + "Delete:"
	log.Info(message + "started")
	res := s.db.Del(sessionId)
	log.Info(message + "ended")
	return res.Err()
}
