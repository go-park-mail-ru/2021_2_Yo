package auth

import (
	protoAuth "backend/microservices/proto/auth"
	"backend/models"
	"context"
)

type AuthService interface {
	Create(ctx context.Context, user models.User) (*models.ResponseBodyUser, error)
	Find(ctx context.Context, user models.User) (*models.ResponseBodyUser, error)
}

type service struct {
	authService protoAuth.AuthClient
}

func NewService (authService protoAuth.AuthClient) AuthService {
	return &service{
		authService: authService,
	}
}

func (s *service) Create(ctx context.Context, user models.User) (*models.ResponseBodyUser, error)  {
	request := protoAuth.User{
		Name: user.Name,
		Surname: user.Surname,
		Mail: user.Mail,
		Password: user.Password,
	}
	userResponse, err := s.authService.CreateUser(ctx,&request) 
	if err != nil {
		return &models.ResponseBodyUser{}, err
	}
	
	response := models.ResponseBodyUser {
		ID: userResponse.ID,
		Name: user.Name,
		Surname: user.Surname,
		Mail: user.Mail,
	}

	return &response, nil
}

func (s *service) Find(ctx context.Context, user models.User) (*models.ResponseBodyUser, error) {
	request := protoAuth.UserAuthData{
		Mail: user.Mail,
		Password: user.Password,
	}
	userResponse, err := s.authService.GetUser(ctx, &request)
	if err != nil {
		return &models.ResponseBodyUser{}, err
	}
	response := models.ResponseBodyUser {
		ID: userResponse.ID,
		Name: user.Name,
		Surname: user.Surname,
		Mail: user.Mail,
		About: user.About,
		ImgUrl: user.ImgUrl,
	}

	return &response, nil
}