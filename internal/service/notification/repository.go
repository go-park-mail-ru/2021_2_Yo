package notification

import "backend/internal/models"

type Repository interface {
	CreateSubscribeNotification(receiverId string, subscriber *models.User, seen bool) error
	CreateInviteNotification(receiverId string, invitor *models.User, event *models.Event, seen bool) error
	CreateNewEventNotification(receiverId string, invitor *models.User, event *models.Event, seen bool) error
	UpdateNotificationsStatus(userId string) error
	GetAllNotifications(userId string) ([]*models.Notification, error)
	GetNewNotifications(userId string) ([]*models.Notification, error)
}
