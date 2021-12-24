package grpc

import (
	proto "backend/internal/microservice/user/proto"
	"backend/internal/models"
	"context"
)

type Repository struct {
	client proto.UserServiceClient
}

func NewRepository(client proto.UserServiceClient) *Repository {
	return &Repository{
		client: client,
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
	if u == nil {
		return nil
	}
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

func (a *Repository) GetUserById(userId string) (*models.User, error) {
	in := &proto.UserId{
		ID: userId,
	}
	out, err := a.client.GetUserById(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := MakeModelUser(out)
	return result, err
}

func (a *Repository) UpdateUserInfo(u *models.User) error {
	in := MakeProtoUser(u)
	_, err := a.client.UpdateUserInfo(context.Background(), in)
	return err
}

func (a *Repository) UpdateUserPassword(userId string, password string) error {
	in := &proto.UpdateUserPasswordRequest{
		ID:       userId,
		Password: password,
	}
	_, err := a.client.UpdateUserPassword(context.Background(), in)
	return err
}

func (a *Repository) GetSubscribers(userId string) ([]*models.User, error) {
	in := &proto.UserId{
		ID: userId,
	}
	out, err := a.client.GetSubscribers(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = MakeModelUser(protoUser)
	}
	return result, err
}

func (a *Repository) GetSubscribes(userId string) ([]*models.User, error) {
	in := &proto.UserId{
		ID: userId,
	}
	out, err := a.client.GetSubscribes(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = MakeModelUser(protoUser)
	}
	return result, err
}

func (a *Repository) GetFriends(userId string, eventId string) ([]*models.User, error) {
	in := &proto.GetFriendsRequest{
		UserId:  userId,
		EventId: eventId,
	}
	out, err := a.client.GetFriends(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = MakeModelUser(protoUser)
	}
	return result, err
}

func (a *Repository) GetVisitors(eventId string) ([]*models.User, error) {
	in := &proto.EventId{
		ID: eventId,
	}
	out, err := a.client.GetVisitors(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = MakeModelUser(protoUser)
	}
	return result, err
}

func (a *Repository) Subscribe(subscribedId string, subscriberId string) error {
	in := &proto.SubscribeRequest{
		SubscribedId: subscribedId,
		SubscriberId: subscriberId,
	}
	_, err := a.client.Subscribe(context.Background(), in)
	return err
}

func (a *Repository) Unsubscribe(subscribedId string, subscriberId string) error {
	in := &proto.SubscribeRequest{
		SubscribedId: subscribedId,
		SubscriberId: subscriberId,
	}
	_, err := a.client.Unsubscribe(context.Background(), in)
	return err
}

func (a *Repository) IsSubscribed(subscribedId string, subscriberId string) (bool, error) {
	in := &proto.SubscribeRequest{
		SubscribedId: subscribedId,
		SubscriberId: subscriberId,
	}
	out, err := a.client.IsSubscribed(context.Background(), in)
	if err != nil {
		return false, err
	}
	result := out.Result
	return result, err
}
