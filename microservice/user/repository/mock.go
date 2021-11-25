package repository

import (
	userGrpc "backend/microservice/user/proto"
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type RepositoryClientMock struct {
	mock.Mock
}

func (m *RepositoryClientMock) GetUserById(ctx context.Context, in *userGrpc.UserId, opts ...grpc.CallOption) (*userGrpc.User, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.User), args.Error(1)
}

func (m *RepositoryClientMock) UpdateUserInfo(ctx context.Context, in *userGrpc.User, opts ...grpc.CallOption) (*userGrpc.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Empty), args.Error(1)
}

func (m *RepositoryClientMock) UpdateUserPassword(ctx context.Context, in *userGrpc.UpdateUserPasswordRequest, opts ...grpc.CallOption) (*userGrpc.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Empty), args.Error(1)
}

func (m *RepositoryClientMock) GetSubscribers(ctx context.Context, in *userGrpc.UserId, opts ...grpc.CallOption) (*userGrpc.Users, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Users), args.Error(1)
}

func (m *RepositoryClientMock) GetSubscribes(ctx context.Context, in *userGrpc.UserId, opts ...grpc.CallOption) (*userGrpc.Users, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Users), args.Error(1)
}

func (m *RepositoryClientMock) GetVisitors(ctx context.Context, in *userGrpc.EventId, opts ...grpc.CallOption) (*userGrpc.Users, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Users), args.Error(1)
}

func (m *RepositoryClientMock) Subscribe(ctx context.Context, in *userGrpc.SubscribeRequest, opts ...grpc.CallOption) (*userGrpc.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Empty), args.Error(1)
}

func (m *RepositoryClientMock) Unsubscribe(ctx context.Context, in *userGrpc.SubscribeRequest, opts ...grpc.CallOption) (*userGrpc.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.Empty), args.Error(1)
}

func (m *RepositoryClientMock) IsSubscribed(ctx context.Context, in *userGrpc.SubscribeRequest, opts ...grpc.CallOption) (*userGrpc.IsSubscribedRequest, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userGrpc.IsSubscribedRequest), args.Error(1)
}
