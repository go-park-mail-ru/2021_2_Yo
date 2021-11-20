package usecase

import (
	proto "backend/microservice/user/proto"
	"backend/models"
	error2 "backend/service/user/error"
	"backend/utils"
	"context"
)

const logMessage = "service:user:usecase:"

type UseCase struct {
	//UserRepositoryClient - это интерфейс, поэтому можно замокать
	userRepo proto.RepositoryClient
	//Потом будет eventRepo
}

func NewUseCase(userRepo proto.RepositoryClient) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

func (a *UseCase) GetUserById(userId string) (*models.User, error) {

	if userId == "" {
		return nil, error2.ErrEmptyData
	}

	in := &proto.UserId{ID: userId}
	protoUser, err := a.userRepo.GetUserById(context.Background(), in)

	if err != nil {
		return nil, err
	}
	resultUser := &models.User{
		ID:       protoUser.ID,
		Name:     protoUser.Name,
		Surname:  protoUser.Surname,
		Mail:     protoUser.Mail,
		Password: protoUser.Password,
		About:    protoUser.About,
		ImgUrl:   protoUser.ImgUrl,
	}
	return resultUser, nil

}

func (a *UseCase) UpdateUserInfo(u *models.User) error {

	if u.ID == "" || u.Name == "" || u.Surname == "" {
		return error2.ErrEmptyData
	}

	in := &proto.User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}

	_, err := a.userRepo.UpdateUserInfo(context.Background(), in)
	return err
}

func (a *UseCase) UpdateUserPassword(userId string, password string) error {
	if userId == "" || password == "" {
		return error2.ErrEmptyData
	}
	hashedPassword := utils.CreatePasswordHash(password)

	in := &proto.UpdateUserPasswordRequest{
		ID:       userId,
		Password: hashedPassword,
	}

	_, err := a.userRepo.UpdateUserPassword(context.Background(), in)
	return err
}

func (a *UseCase) Subscribe(subscribedId string, subscriberId string) error {

	if subscribedId == "" || subscriberId == "" {
		return error2.ErrEmptyData
	}

	in := &proto.SubscribeRequest{
		SubscribedId: subscribedId,
		SubscriberId: subscriberId,
	}

	_, err := a.userRepo.Subscribe(context.Background(), in)

	return err

}

func (a *UseCase) GetSubscribers(userId string) ([]*models.User, error) {

	if userId == "" {
		return nil, error2.ErrEmptyData
	}

	in := &proto.UserId{
		ID: userId,
	}

	out, err := a.userRepo.GetSubscribers(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = &models.User{
			ID:       protoUser.ID,
			Name:     protoUser.Name,
			Surname:  protoUser.Surname,
			Mail:     protoUser.Mail,
			Password: protoUser.Password,
			About:    protoUser.About,
			ImgUrl:   protoUser.ImgUrl,
		}
	}

	return result, nil

}

func (a *UseCase) GetSubscribes(userId string) ([]*models.User, error) {

	if userId == "" {
		return nil, error2.ErrEmptyData
	}

	in := &proto.UserId{
		ID: userId,
	}

	out, err := a.userRepo.GetSubscribes(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.User, len(out.Users))
	for i, protoUser := range out.Users {
		result[i] = &models.User{
			ID:       protoUser.ID,
			Name:     protoUser.Name,
			Surname:  protoUser.Surname,
			Mail:     protoUser.Mail,
			Password: protoUser.Password,
			About:    protoUser.About,
			ImgUrl:   protoUser.ImgUrl,
		}
	}

	return result, nil

}
