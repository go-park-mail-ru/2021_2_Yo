package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"context"
	"github.com/stretchr/testify/mock"
)

type AuthClientMock struct {
	mock.Mock
}

func (m *AuthClientMock) SignUp(ctx context.Context, in *protoAuth.SignUpRequest) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) SignIn(ctx context.Context, in *protoAuth.SignInRequest) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) CreateSession(ctx context.Context, in *protoAuth.UserId) (*protoAuth.Session, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.Session), args.Error(1)
}

func (m *AuthClientMock) CheckSession(ctx context.Context, in *protoAuth.Session) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) DeleteSession(ctx context.Context, in *protoAuth.Session) (*protoAuth.Success, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.Success), args.Error(1)
}

func (m *AuthClientMock) CreateToken(ctx context.Context, in *protoAuth.UserId) (*protoAuth.CSRFToken, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.CSRFToken), args.Error(1)
}

func (m *AuthClientMock) CheckToken(ctx context.Context, in *protoAuth.CSRFToken) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}
