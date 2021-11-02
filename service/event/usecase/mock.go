package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func (m *UseCaseMock) List() ([]*models.Event, error) {
	return nil, nil
}
