package usecase

import (
	"backend/models"
	"github.com/stretchr/testify/mock"
)

type UseCaseEventsManagerMock struct {
	mock.Mock
}

func (m *UseCaseEventsManagerMock) List() ([]*models.Event, error) {
	return nil, nil
}
