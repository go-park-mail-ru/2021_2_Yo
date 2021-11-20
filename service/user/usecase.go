package user

import "backend/models"

type UseCase interface {
	GetUserById(userId string) (*models.User, error)
	UpdateUserInfo(user *models.User) error
	UpdateUserPassword(userId string, password string) error
	Subscribe(subscribedId string, subscriberId string) error
}
