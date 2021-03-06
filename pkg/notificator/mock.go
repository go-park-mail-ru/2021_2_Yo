package notificator

import (
	"backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type NotificatorMock struct {
	mock.Mock
}

func (m *NotificatorMock) NewSubscriberNotification(receiverId string, userId string) error {
	args := m.Called(receiverId, userId)
	return args.Error(0)
}

func (m *NotificatorMock) DeleteSubscribeNotification(receiverId string, userId string) error {
	args := m.Called(receiverId, userId)
	return args.Error(0)
}

func (m *NotificatorMock) InvitationNotification(receiverId string, userId string, eventId string) error {
	args := m.Called(receiverId, userId, eventId)
	return args.Error(0)
}

func (m *NotificatorMock) NewEventNotification(userId string, eventId string) error {
	args := m.Called(userId, eventId)
	return args.Error(0)
}

func (m *NotificatorMock) UpdateNotificationsStatus(receiverId string) error {
	args := m.Called(receiverId)
	return args.Error(0)
}

func (m *NotificatorMock) GetAllNotifications(receiverId string) ([]*models.Notification, error) {
	args := m.Called(receiverId)
	return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *NotificatorMock) GetNewNotifications(receiverId string) ([]*models.Notification, error) {
	args := m.Called(receiverId)
	return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *NotificatorMock) EventTomorrowNotification() error {
	args := m.Called()
	return args.Error(0)
}

func (m *NotificatorMock) PingConnections() int {
	args := m.Called()
	return args.Get(0).(int)
}
