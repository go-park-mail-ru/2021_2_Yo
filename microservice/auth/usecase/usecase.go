package usecase

import (
	"backend/microservice/auth/interfaces"
	protoAuth "backend/microservice/auth/proto"
	"backend/pkg/models"
	"backend/pkg/utils"
	"context"

	log "github.com/sirupsen/logrus"
)

const (
	logMessage = "microservice auth"
)

type authService struct {
	authUserRepository    interfaces.UserRepository
	authSessionRepository interfaces.SessionRepository
}

func NewService(authUserRepository interfaces.UserRepository, authSessionRepository interfaces.SessionRepository) *authService {
	return &authService{
		authUserRepository:    authUserRepository,
		authSessionRepository: authSessionRepository,
	}
}

func (s *authService) SignUp(ctx context.Context, in *protoAuth.SignUpRequest) (*protoAuth.UserId, error) {
	message := logMessage + "SignUp:"
	log.Debug(message + "started")

	log.Debug(message+"in = ", in)

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
	message := logMessage + "SignIn:"
	log.Debug(message + "started")

	log.Debug(message+"in = ", in)

	u, err := s.authUserRepository.GetUser(in.Mail, utils.CreatePasswordHash(in.Password))
	if err != nil {
		return &protoAuth.UserId{}, err
	}

	out := &protoAuth.UserId{
		ID: u.ID,
	}
	return out, nil
}
