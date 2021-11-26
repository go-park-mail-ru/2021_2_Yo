package client

import (
	proto "backend/microservice/user/proto"
	"backend/pkg/models"
	"backend/service/user"
	"context"
)

type UserService struct {
	repository user.Repository
}

func NewUserService(repository user.Repository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func MakeProtoUser(u *models.User) *proto.User {
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

func MakeModelUser(u *proto.User) *models.User {
	return &models.User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}
}

func MakeProtoUsers(u []*models.User) *proto.Users {
	if u == nil {
		return &proto.Users{}
	}
	result := make([]*proto.User, len(u))
	for i, modelUser := range u {
		result[i] = MakeProtoUser(modelUser)
	}
	return &proto.Users{
		Users: result,
	}
}

func (c *UserService) GetUserById(ctx context.Context, in *proto.UserId) (*proto.User, error) {
	userId := in.ID
	modelUser, err := c.repository.GetUserById(userId)
	out := MakeProtoUser(modelUser)
	return out, err
}

func (c *UserService) UpdateUserInfo(ctx context.Context, in *proto.User) (*proto.Empty, error) {
	modelUser := MakeModelUser(in)
	err := c.repository.UpdateUserInfo(modelUser)
	out := &proto.Empty{}
	return out, err
}

func (c *UserService) UpdateUserPassword(ctx context.Context, in *proto.UpdateUserPasswordRequest) (*proto.Empty, error) {
	userId := in.ID
	password := in.Password
	err := c.repository.UpdateUserPassword(userId, password)
	out := &proto.Empty{}
	return out, err
}

func (c *UserService) GetSubscribers(ctx context.Context, in *proto.UserId) (*proto.Users, error) {
	userId := in.ID
	modelUsers, err := c.repository.GetSubscribers(userId)
	out := MakeProtoUsers(modelUsers)
	return out, err
}

func (c *UserService) GetSubscribes(ctx context.Context, in *proto.UserId) (*proto.Users, error) {
	userId := in.ID
	modelUsers, err := c.repository.GetSubscribes(userId)
	out := MakeProtoUsers(modelUsers)
	return out, err
}

func (c *UserService) GetVisitors(ctx context.Context, in *proto.EventId) (*proto.Users, error) {
	eventId := in.ID
	modelUsers, err := c.repository.GetVisitors(eventId)
	out := MakeProtoUsers(modelUsers)
	return out, err
}

func (c *UserService) Subscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.Empty, error) {
	subscribedId := in.SubscribedId
	subscriberId := in.SubscriberId
	err := c.repository.Subscribe(subscribedId, subscriberId)
	out := &proto.Empty{}
	return out, err
}

func (c *UserService) Unsubscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.Empty, error) {
	subscribedId := in.SubscribedId
	subscriberId := in.SubscriberId
	err := c.repository.Unsubscribe(subscribedId, subscriberId)
	out := &proto.Empty{}
	return out, err
}

func (c *UserService) IsSubscribed(ctx context.Context, in *proto.SubscribeRequest) (*proto.IsSubscribedRequest, error) {
	subscribedId := in.SubscribedId
	subscriberId := in.SubscriberId
	result, err := c.repository.IsSubscribed(subscribedId, subscriberId)
	out := &proto.IsSubscribedRequest{
		Result: result,
	}
	return out, err
}
