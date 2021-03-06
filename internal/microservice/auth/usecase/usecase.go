package usecase

import (
	interfaces2 "backend/internal/microservice/auth/interfaces"
	protoAuth "backend/internal/microservice/auth/proto"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
)

//const logMessage = "microservice auth"

type authService struct {
	authUserRepository    interfaces2.UserRepository
	authSessionRepository interfaces2.SessionRepository
}

func NewService(authUserRepository interfaces2.UserRepository, authSessionRepository interfaces2.SessionRepository) *authService {
	return &authService{
		authUserRepository:    authUserRepository,
		authSessionRepository: authSessionRepository,
	}
}

func (s *authService) SignUp(ctx context.Context, in *protoAuth.SignUpRequest) (*protoAuth.UserId, error) {

	newUser := models.User{
		Name:     in.Name,
		Surname:  in.Surname,
		Mail:     in.Mail,
		Password: utils.CreatePasswordHash(in.Password),
	}

	userId, err := s.authUserRepository.CreateUser(&newUser)
	if err != nil {
		return &protoAuth.UserId{}, err
	}

	out := &protoAuth.UserId{ID: userId}
	return out, nil
}

func (s *authService) SignIn(ctx context.Context, in *protoAuth.SignInRequest) (*protoAuth.UserId, error) {

	u, err := s.authUserRepository.GetUser(in.Mail, utils.CreatePasswordHash(in.Password))
	if err != nil {
		return &protoAuth.UserId{}, err
	}

	out := &protoAuth.UserId{
		ID: u.ID,
	}
	return out, nil
}
