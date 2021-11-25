package repository

import (
	proto "backend/microservice/user/proto"
	"backend/pkg/models"
	error2 "backend/service/user/error"
	"strconv"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Mail     string `db:"mail"`
	Password string `db:"password"`
	About    string `db:"about"`
	ImgUrl   string `db:"img_url"`
}

func toPostgresUser(u *models.User) (*User, error) {
	var userIdInt int
	if u.ID == "" {
		userIdInt = 0
	} else {
		tempUserId, err := strconv.Atoi(u.ID)
		if err != nil {
			return nil, error2.ErrAtoi
		}
		userIdInt = tempUserId
	}
	return &User{
		ID:       userIdInt,
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}, nil
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       strconv.Itoa(u.ID),
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}
}

func toProtoUser(u *models.User) *proto.User {
	return &proto.User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}
}
