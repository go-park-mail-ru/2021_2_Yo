package session

import (
	authServiceModels "backend/microservice/auth/models"
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
	res := s.db.Set(data.SessionId, data.UserId, data.Expiration)
	return res.Err()
}

func (s *Repository) Check(sessionId string) (string, error) {
	res := s.db.Get(sessionId)
	return res.Val(), res.Err()
}

func (s *Repository) Delete(sessionId string) error {
	res := s.db.Del(sessionId)
	return res.Err()
}
