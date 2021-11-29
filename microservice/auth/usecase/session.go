package usecase

import (
	authServiceModels "backend/microservice/auth/models"
	protoAuth "backend/microservice/auth/proto"
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrEmptySessionId = errors.New("session id is empty")
	letterRunes       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func generateSessionId(n int) string {
	if n == 0 {
		return ""
	}
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const (
	sessionIdLength = 16
	sessionLifeTime = time.Hour * 24
)

func (s *authService) CreateSession(ctx context.Context, protoUserId *protoAuth.UserId) (*protoAuth.Session, error) {
	sessionData := &authServiceModels.SessionData{
		SessionId:  generateSessionId(sessionIdLength),
		UserId:     protoUserId.ID,
		Expiration: sessionLifeTime,
	}
	id, _ := strconv.Atoi(sessionData.UserId)
	if id <= 0 {
		sessionData.SessionId = ""
	}
	err := s.authSessionRepository.Create(sessionData)
	if err != nil {
		return &protoAuth.Session{}, err
	}
	response := &protoAuth.Session{
		Session: sessionData.SessionId,
	}
	return response, err
}

func (s *authService) CheckSession(ctx context.Context, protoSession *protoAuth.Session) (*protoAuth.UserId, error) {
	sessionId := protoSession.Session
	if sessionId == "" {
		return &protoAuth.UserId{}, ErrEmptySessionId
	}
	userId, err := s.authSessionRepository.Check(sessionId)
	if err != nil {
		return &protoAuth.UserId{}, err
	}
	response := &protoAuth.UserId{
		ID: userId,
	}
	return response, nil
}

func (s *authService) DeleteSession(ctx context.Context, protoSession *protoAuth.Session) (*protoAuth.Success, error) {
	err := s.authSessionRepository.Delete(protoSession.Session)
	if err != nil {
		return &protoAuth.Success{}, err
	}
	response := &protoAuth.Success{
		Ok: "success",
	}
	return response, nil
}
