package csrf

import (
	my_error "backend/csrf/error"
	"backend/csrf/models"
	"backend/csrf/repository"
	"time"

	"github.com/satori/go.uuid"
)
type Manager struct {
	repository repository.Repository
}

func NewManager(repo repository.Repository) *Manager {
	return &Manager{
		repository: repo,
	}
}

const sessionLifeTime = time.Hour * 24 * 7 //Неделя

func generateCSRFToken() string {
	csrfToken := uuid.NewV4()
	return csrfToken.String()
}

func (m *Manager) Create(userID string) (string, error) {
	CSRFData := &models.CSRFData{
		CSRFToken:  generateCSRFToken(),
		UserId:     userID,
		Expiration: int(sessionLifeTime),
	}

	err := m.repository.Create(CSRFData)
	if err != nil {
		return "", err
	}
	return CSRFData.CSRFToken, nil
}

func (m *Manager) Check(susToken string) (string, error) {
	if susToken == "" {
		return "", my_error.ErrEmptyToken
	}
	userId, err := m.repository.Check(susToken)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (m *Manager) Delete(susToken string) error {
	err := m.repository.Delete(susToken)
	if err != nil {
		return err
	}
	return nil
}