package eventRepository

import (
	proto "backend/microservice/event/proto"
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type RepositoryClientMock struct {
	mock.Mock
}

func (m *RepositoryClientMock) CreateEvent(ctx context.Context, in *proto.Event, opts ...grpc.CallOption) (*proto.EventId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.EventId), args.Error(1)
}

func (m *RepositoryClientMock) UpdateEvent(ctx context.Context, in *proto.UpdateEventRequest, opts ...grpc.CallOption) (*proto.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Empty), args.Error(1)
}

func (m *RepositoryClientMock) DeleteEvent(ctx context.Context, in *proto.DeleteEventRequest, opts ...grpc.CallOption) (*proto.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Empty), args.Error(1)
}

func (m *RepositoryClientMock) GetEventById(ctx context.Context, in *proto.EventId, opts ...grpc.CallOption) (*proto.Event, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Event), args.Error(1)
}

func (m *RepositoryClientMock) GetEvents(ctx context.Context, in *proto.GetEventsRequest, opts ...grpc.CallOption) (*proto.Events, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Events), args.Error(1)
}

func (m *RepositoryClientMock) GetVisitedEvents(ctx context.Context, in *proto.UserId, opts ...grpc.CallOption) (*proto.Events, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Events), args.Error(1)
}

func (m *RepositoryClientMock) GetCreatedEvents(ctx context.Context, in *proto.UserId, opts ...grpc.CallOption) (*proto.Events, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Events), args.Error(1)
}

func (m *RepositoryClientMock) Visit(ctx context.Context, in *proto.VisitRequest, opts ...grpc.CallOption) (*proto.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Empty), args.Error(1)
}

func (m *RepositoryClientMock) Unvisit(ctx context.Context, in *proto.VisitRequest, opts ...grpc.CallOption) (*proto.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.Empty), args.Error(1)
}

func (m *RepositoryClientMock) IsVisited(ctx context.Context, in *proto.VisitRequest, opts ...grpc.CallOption) (*proto.IsVisitedRequest, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.IsVisitedRequest), args.Error(1)
}

func (m *RepositoryClientMock) GetCities(ctx context.Context, in *proto.Empty, opts ...grpc.CallOption) (*proto.GetCitiesRequest, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.GetCitiesRequest), args.Error(1)
}
