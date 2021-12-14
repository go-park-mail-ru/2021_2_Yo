package session

import (
	authServiceModels "backend/internal/microservice/auth/models"
	"github.com/go-redis/redis"
	log "backend/pkg/logger"
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
	log.Info(message)
	res := s.db.Set(data.SessionId, data.UserId, data.Expiration)
	return res.Err()
}

func (s *Repository) Check(sessionId string) (string, error) {
	message := logMessage + "Check:"
	log.Info(message)
	res := s.db.Get(sessionId)
	return res.Val(), res.Err()
}

func (s *Repository) Delete(sessionId string) error {
	message := logMessage + "Delete:"
	log.Info(message)
	res := s.db.Del(sessionId)
	return res.Err()
}
