package usecase

import (
	authServiceModels "backend/microservice/auth/models"
	"github.com/stretchr/testify/mock"
)

type AuthSessionMock struct {
	mock.Mock
}

func (m *AuthSessionMock) Create(data *authServiceModels.SessionData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *AuthSessionMock) Check(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.String(0), args.Error(1)
}

func (m *AuthSessionMock) Delete(sessionId string) error {
	args := m.Called(sessionId)
	return args.Error(0)
}

/*
func (m *AuthClientMock) CreateToken(ctx context.Context, in *protoAuth.UserId, opts ...grpc.CallOption) (*protoAuth.CSRFToken, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.CSRFToken), args.Error(1)
}

func (m *AuthClientMock) CheckToken(ctx context.Context, in *protoAuth.CSRFToken, opts ...grpc.CallOption) (*protoAuth.UserId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoAuth.UserId), args.Error(1)
}
*/
