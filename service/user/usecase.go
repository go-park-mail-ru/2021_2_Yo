package user

import (
	"backend/pkg/models"
)

type UseCase interface {
	GetUserById(userId string) (*models.User, error)
	UpdateUserInfo(user *models.User) error
	UpdateUserPassword(userId string, password string) error
	Subscribe(subscribedId string, subscriberId string) error
	GetSubscribers(userId string) ([]*models.User, error)
	GetSubscribes(userId string) ([]*models.User, error)
	GetVisitors(eventId string) ([]*models.User, error)
}
