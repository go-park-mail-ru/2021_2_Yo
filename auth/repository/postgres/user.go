package postgres

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
	About    string
}

func toPostgresUser(u *models.User) *User {
	return &User{
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
	}
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       strconv.Itoa(u.ID),
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
	}
}
