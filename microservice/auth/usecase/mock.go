package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type AuthClientMock struct {
	mock.Mock
}

func (m *AuthClientMock) SignUp(ctx context.Context, in *protoAuth.SignUpRequest, opts ...grpc.CallOption) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) SignIn(ctx context.Context, in *protoAuth.SignInRequest, opts ...grpc.CallOption) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) CreateSession(ctx context.Context, in *protoAuth.UserId, opts ...grpc.CallOption) (*protoAuth.Session, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.Session), args.Error(1)
}

func (m *AuthClientMock) CheckSession(ctx context.Context, in *protoAuth.Session, opts ...grpc.CallOption) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}

func (m *AuthClientMock) DeleteSession(ctx context.Context, in *protoAuth.Session, opts ...grpc.CallOption) (*protoAuth.Success, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.Success), args.Error(1)
}

func (m *AuthClientMock) CreateToken(ctx context.Context, in *protoAuth.UserId, opts ...grpc.CallOption) (*protoAuth.CSRFToken, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.CSRFToken), args.Error(1)
}

func (m *AuthClientMock) CheckToken(ctx context.Context, in *protoAuth.CSRFToken, opts ...grpc.CallOption) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}