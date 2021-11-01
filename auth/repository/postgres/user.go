package postgres

import (
	"backend/models"
	"strconv"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Mail     string `db:"mail"`
	Password string `db:"password"`
	About    string `db:"about"`
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
