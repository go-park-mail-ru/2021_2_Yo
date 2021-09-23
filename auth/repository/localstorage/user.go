package localstorage

import (
	"backend/models"
	"strconv"
)

type User struct {
	ID       int
	Username string
	Password string
}

func toLocalstorageUser(u *models.User) *User {
	return &User{
		Username: u.Username,
		Password: u.Password,
	}
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       strconv.Itoa(u.ID),
		Username: u.Username,
		Password: u.Password,
	}
}
