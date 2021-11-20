package auth

import (
	protoAuth "backend/microservice/auth/proto"
	models2 "backend/pkg/models"
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, user models2.User) (*models2.UserResponseBody, error)
	SignIn(ctx context.Context, user models2.User) (string, error)
	CreateSession(ctx context.Context, userId string) (string, error)
	CheckSession(ctx context.Context, sessionId string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
	CreateToken(ctx context.Context, userId string) (string, error)
	CheckToken(ctx context.Context, csrfToken string) (string, error)
}

type service struct {
	authService protoAuth.AuthClient
}

func NewService(authService protoAuth.AuthClient) AuthService {
	return &service{
		authService: authService,
	}
}

func (s *service) SignUp(ctx context.Context, user models2.User) (*models2.UserResponseBody, error) {
	request := protoAuth.User{
		Name:     user.Name,
		Surname:  user.Surname,
		Mail:     user.Mail,
		Password: user.Password,
	}
	userResponse, err := s.authService.CreateUser(ctx, &request)
	if err != nil {
		return &models2.UserResponseBody{}, err
	}

	response := models2.UserResponseBody{
		ID:      userResponse.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Mail:    user.Mail,
	}

	return &response, nil
}

func (s *service) SignIn(ctx context.Context, user models2.User) (string, error) {
	request := protoAuth.UserAuthData{
		Mail:     user.Mail,
		Password: user.Password,
	}
	userResponse, err := s.authService.GetUser(ctx, &request)
	if err != nil {
		return "", err
	}
	response := userResponse.UserId

	return response, nil
}

func (s *service) CreateSession(ctx context.Context, userId string) (string, error) {
	request := protoAuth.UserId{
		UserId: userId,
	}
	sessionResponse, err := s.authService.CreateSession(ctx, &request)
	if err != nil {
		return "", err
	}
	response := sessionResponse.Session
	return response, nil
}

func (s *service) CheckSession(ctx context.Context, SessionId string) (string, error) {
	request := protoAuth.Session{
		Session: SessionId,
	}
	userIdResponse, err := s.authService.CheckSession(ctx, &request)
	if err != nil {
		return "", err
	}
	response := userIdResponse.UserId
	return response, nil
}

func (s *service) DeleteSession(ctx context.Context, SessionId string) error {
	request := protoAuth.Session{
		Session: SessionId,
	}
	_, err := s.authService.DeleteSession(ctx, &request)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateToken(ctx context.Context, userId string) (string, error) {
	request := protoAuth.UserId{
		UserId: userId,
	}
	tokenResponse, err := s.authService.CreateToken(ctx, &request)
	if err != nil {
		return "", err
	}
	response := tokenResponse.CSRFToken
	return response, nil
}

func (s *service) CheckToken(ctx context.Context, csrfToken string) (string, error) {
	request := protoAuth.CSRFToken{
		CSRFToken: csrfToken,
	}
	userIdResponse, err := s.authService.CheckToken(ctx, &request)
	if err != nil {
		return "", err
	}
	response := userIdResponse.UserId
	return response, nil
}
