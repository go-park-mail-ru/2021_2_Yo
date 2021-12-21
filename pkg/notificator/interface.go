package notificator

import "backend/internal/models"

type NotificationManager interface {
	NewSubscriberNotification(receiverId string, userId string) error
	DeleteSubscribeNotification(receiverId string, userId string) error
	InvitationNotification(receiverId string, userId string, eventId string) error
	NewEventNotification(userId string, eventId string) error
	UpdateNotificationsStatus(receiverId string) error
	GetAllNotifications(receiverId string) ([]*models.Notification, error)
	GetNewNotifications(receiverId string) ([]*models.Notification, error)
}
