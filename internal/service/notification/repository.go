package notification

import (
	"backend/internal/models"
)

type Repository interface {
	CreateSubscribeNotification(receiverId string, subscriber *models.User, event *models.Event) error
	DeleteSubscribeNotification(receiverId string, userId string) error
	CreateInviteNotification(receiverId string, invitor *models.User, event *models.Event) error
	CreateNewEventNotification(receiverId string, invitor *models.User, event *models.Event) error
	UpdateNotificationsStatus(userId string) error
	GetAllNotifications(userId string) ([]*models.Notification, error)
	GetNewNotifications(userId string) ([]*models.Notification, error)
	CreateTomorrowEventNotification(receiverId string, invitor *models.User, event *models.Event) error
}
