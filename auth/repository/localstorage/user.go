package localstorage

import (
	"backend/models"
	"strconv"
)

type User struct {
	ID       int
	Name     string
	Surname  string
	Mail     string
	Password string
}

func toLocalstorageUser(u *models.User) *User {
	return &User{
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
	}
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       strconv.Itoa(u.ID),
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
	}
}
