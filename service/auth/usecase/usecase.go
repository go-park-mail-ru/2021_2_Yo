package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"backend/pkg/models"
	"context"
	log "backend/pkg/logger"
)

const logMessage = "service:auth:usecase:"

type UseCase struct {
	client protoAuth.AuthClient
}

func NewUseCase(authService protoAuth.AuthClient) *UseCase {
	return &UseCase{
		client: authService,
	}
}

func (s *UseCase) SignUp(u *models.User) (string, error) {
	log.Debug(logMessage+"SignUp:started")
	in := &protoAuth.SignUpRequest{
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
	}
	out, err := s.client.SignUp(context.Background(), in)
	if err != nil {
		return "", err
	}
	userId := out.ID
	log.Debug(logMessage+"SignUp:ended")
	return userId, nil
}

func (s *UseCase) SignIn(u *models.User) (string, error) {
	in := &protoAuth.SignInRequest{
		Mail:     u.Mail,
		Password: u.Password,
	}
	out, err := s.client.SignIn(context.Background(), in)
	if err != nil {
		return "", err
	}
	userId := out.ID
	return userId, nil
}

func (s *UseCase) CreateSession(userId string) (string, error) {
	in := &protoAuth.UserId{
		ID: userId,
	}
	out, err := s.client.CreateSession(context.Background(), in)
	if err != nil {
		return "", err
	}
	sessionId := out.Session
	return sessionId, nil
}

func (s *UseCase) CheckSession(SessionId string) (string, error) {
	in := &protoAuth.Session{
		Session: SessionId,
	}
	out, err := s.client.CheckSession(context.Background(), in)
	if err != nil {
		return "", err
	}
	userId := out.ID
	return userId, nil
}

func (s *UseCase) DeleteSession(SessionId string) error {
	in := &protoAuth.Session{
		Session: SessionId,
	}
	_, err := s.client.DeleteSession(context.Background(), in)
	if err != nil {
		return err
	}
	return nil
}

func (s *UseCase) CreateToken(userId string) (string, error) {
	in := &protoAuth.UserId{
		ID: userId,
	}
	out, err := s.client.CreateToken(context.Background(), in)
	if err != nil {
		return "", err
	}
	token := out.CSRFToken
	return token, nil
}

func (s *UseCase) CheckToken(csrfToken string) (string, error) {
	in := &protoAuth.CSRFToken{
		CSRFToken: csrfToken,
	}
	out, err := s.client.CheckToken(context.Background(), in)
	if err != nil {
		return "", err
	}
	userId := out.ID
	return userId, nil
}
