package manager

import (
	log "backend/logger"
	error2 "backend/service/session/error"
	"backend/service/session/models"
	"backend/service/session/repository"
	"math/rand"
	"time"
)

type Manager struct {
	repository repository.Repository
}

func NewManager(repo repository.Repository) *Manager {
	return &Manager{
		repository: repo,
	}
}

const (
	logMessage      = "service:session:manager:"
	sessionIdLength = 16
	sessionLifeTime = time.Hour * 24
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateSessionId(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (m *Manager) Create(userId string) (string, error) {
	sessionData := &models.SessionData{
		SessionId:  generateSessionId(sessionIdLength),
		UserId:     userId,
		Expiration: sessionLifeTime,
	}
	err := m.repository.Create(sessionData)
	if err != nil {
		return "", err
	}
	return sessionData.SessionId, nil
}

func (m *Manager) Check(sessionId string) (string, error) {
	if sessionId == "" {
		return "", error2.ErrEmptySessionId
	}
	userId, err := m.repository.Check(sessionId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (m *Manager) Delete(sessionId string) error {
	err := m.repository.Delete(sessionId)
	if err != nil {
		log.Error(logMessage+"Delete:err =", err)
		return err
	}
	return nil
}
