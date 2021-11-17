package usecase

import (
	"backend/microservices/auth/interfaces"
	protoAuth "backend/microservices/proto/auth"
	"backend/models"
	"backend/utils"
	"context"

	log "github.com/sirupsen/logrus"
)

const (
	logMessage = "microservice auth"
)

type authService struct {
	authUserRepository interfaces.UserRepository
	authSessionRepository interfaces.SessionRepository
}

func NewService(authUserRepository interfaces.UserRepository, authSessionRepository interfaces.SessionRepository) *authService {
	return &authService {
		authUserRepository: authUserRepository,
		authSessionRepository: authSessionRepository,
	}
}



func userConvertToProto(user *models.User) protoAuth.SuccessUserResponse {
	return protoAuth.SuccessUserResponse {
		ID: user.ID,
		Name: user.Name,
		Surname: user.Surname,
		Mail: user.Mail,
		Password: user.Password,
		About: user.About,

	}
}

func (s *authService) CreateUser(ctx context.Context, protoUser *protoAuth.User) (*protoAuth.SuccessUserResponse, error) {
	message := logMessage + "SignUp:"
	log.Debug(message + "started")
	log.Info(protoUser)
	newUser := models.User {
		Name: protoUser.Name,
		Surname: protoUser.Surname,
		Mail: protoUser.Mail,
		Password: utils.CreatePasswordHash(protoUser.Password),
	}

	var err error
	newUser.ID, err = s.authUserRepository.CreateUser(&newUser)

	if err != nil {
		return &protoAuth.SuccessUserResponse{}, err
	}
	response := userConvertToProto(&newUser)
	return &response, nil
}

func (s *authService) GetUser(ctx context.Context, protoUser *protoAuth.UserAuthData) (*protoAuth.UserId, error) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	user, err := s.authUserRepository.GetUser(protoUser.Mail,utils.CreatePasswordHash(protoUser.Password))

	if err != nil {
		return &protoAuth.UserId{}, err
	}
	response := &protoAuth.UserId{
		UserId: user.ID,
	}
	return response, nil
}

