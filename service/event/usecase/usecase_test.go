package usecase

import (
	eventGrpc "backend/microservice/event/proto"
	repository "backend/microservice/event/repository"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:event:usecase:"

var createEventTests = []struct {
	id            int
	event         *models.Event
	outputErr     error
	outputEventId string
}{
	{1,
		&models.Event{
			AuthorId: "test",
			Geo:      "(1.23232323, 4.3223232323)",
			Tag:      []string{"test"},
		},
		nil,
		"",
	},
	{2,
		nil,
		error2.ErrEmptyData,
		"",
	},
	{3,
		&models.Event{
			AuthorId: "test",
			Geo:      "(1.23232323, 4.3223232323)",
		},
		errors.New("test_err"),
		"",
	},
}

func TestCreateEvent(t *testing.T) {
	for _, test := range createEventTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		var in *eventGrpc.Event
		if test.event != nil {
			in = MakeProtoEvent(test.event)
		}
		repositoryMock.On("CreateEvent", context.Background(), in).Return(&eventGrpc.EventId{}, test.outputErr)
		actualEventId, actualErr := useCaseTest.CreateEvent(test.event)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputEventId, actualEventId, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateEventTests = []struct {
	id        int
	event     *models.Event
	userId    string
	outputErr error
}{
	{1,
		&models.Event{
			ID:       "test",
			AuthorId: "test",
			Geo:      "(1.23232323, 4.3223232323)",
			Tag:      []string{"test"},
		},
		"test",
		nil,
	},
	{2,
		nil,
		"",
		error2.ErrEmptyData,
	},
	{3,
		&models.Event{
			ID:       "test",
			AuthorId: "test",
			Geo:      "(1.23232323, 4.3223232323)",
		},
		"test",
		errors.New("test_err"),
	},
}

func TestUpdateEvent(t *testing.T) {
	for _, test := range updateEventTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		var pe *eventGrpc.Event
		if test.event != nil {
			pe = MakeProtoEvent(test.event)
		}
		in := &eventGrpc.UpdateEventRequest{
			Event:  pe,
			UserId: test.userId,
		}
		repositoryMock.On("UpdateEvent", context.Background(), in).Return(&eventGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.UpdateEvent(test.event, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var deleteEventTests = []struct {
	id        int
	eventId   string
	userId    string
	outputErr error
}{
	{1,
		"test",
		"test",
		nil,
	},
	{2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestDeleteEvent(t *testing.T) {
	for _, test := range deleteEventTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.DeleteEventRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		repositoryMock.On("DeleteEvent", context.Background(), in).Return(&eventGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.DeleteEvent(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getEventByIdTests = []struct {
	id        int
	eventId   string
	outputErr error
	outputRes *models.Event
}{
	{1,
		"test",
		nil,
		&models.Event{},
	},
	{2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetEventById(t *testing.T) {
	for _, test := range getEventByIdTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.EventId{
			ID: test.eventId,
		}
		repositoryMock.On("GetEventById", context.Background(), in).Return(&eventGrpc.Event{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetEventById(test.eventId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getEventsTests = []struct {
	id        int
	title     string
	category  string
	city      string
	date      string
	tags      []string
	outputErr error
	outputRes []*models.Event
}{
	{1,
		"test",
		"test",
		"test",
		"test",
		[]string{"test"},
		nil,
		[]*models.Event{},
	},
	{1,
		"test",
		"test",
		"test",
		"test",
		nil,
		errors.New("test_err"),
		nil,
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.GetEventsRequest{
			Title:    test.title,
			Category: test.category,
			City:     test.city,
			Date:     test.date,
			Tags:     test.tags,
		}
		repositoryMock.On("GetEvents", context.Background(), in).Return(&eventGrpc.Events{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetEvents(test.title, test.category, test.city, test.date, test.tags)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getVisitedEventsTests = []struct {
	id        int
	userId    string
	outputErr error
	outputRes []*models.Event
}{
	{1,
		"test",
		nil,
		[]*models.Event{
			&models.Event{},
		},
	},
	{2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetVisitedEvents(t *testing.T) {
	for _, test := range getVisitedEventsTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.UserId{
			ID: test.userId,
		}
		repositoryMock.On("GetVisitedEvents", context.Background(), in).Return(&eventGrpc.Events{
			Events: []*eventGrpc.Event{
				&eventGrpc.Event{},
			},
		}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetVisitedEvents(test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getCreatedEventsTests = []struct {
	id        int
	userId    string
	outputErr error
	outputRes []*models.Event
}{
	{1,
		"test",
		nil,
		[]*models.Event{
			&models.Event{},
		},
	},
	{2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetCreatedEvents(t *testing.T) {
	for _, test := range getCreatedEventsTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.UserId{
			ID: test.userId,
		}
		repositoryMock.On("GetCreatedEvents", context.Background(), in).Return(&eventGrpc.Events{
			Events: []*eventGrpc.Event{
				&eventGrpc.Event{},
			},
		}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetCreatedEvents(test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes)
	}
}

var visitTests = []struct {
	id        int
	userId    string
	eventId   string
	outputErr error
}{
	{1,
		"test",
		"test",
		nil,
	},
	{2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestVisit(t *testing.T) {
	for _, test := range visitTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		repositoryMock.On("Visit", context.Background(), in).Return(&eventGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.Visit(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var unvisitTests = []struct {
	id        int
	userId    string
	eventId   string
	outputErr error
}{
	{1,
		"test",
		"test",
		nil,
	},
	{2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestUnvisit(t *testing.T) {
	for _, test := range unvisitTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		repositoryMock.On("Unvisit", context.Background(), in).Return(&eventGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.Unvisit(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var isVisitedTests = []struct {
	id        int
	userId    string
	eventId   string
	outputErr error
	outputRes bool
}{
	{1,
		"test",
		"test",
		nil,
		true,
	},
	{2,
		"",
		"",
		error2.ErrEmptyData,
		false,
	},
	{3,
		"test",
		"test",
		errors.New("test_err"),
		false,
	},
}

func TestIsVisited(t *testing.T) {
	for _, test := range isVisitedTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		repositoryMock.On("IsVisited", context.Background(), in).Return(&eventGrpc.IsVisitedRequest{
			Result: test.outputRes,
		}, test.outputErr)
		actualRes, actualErr := useCaseTest.IsVisited(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getCitiesTests = []struct {
	id        int
	outputErr error
	outputRes []string
}{
	{1,
		nil,
		[]string{"test"},
	},
	{2,
		error2.ErrEmptyData,
		nil,
	},
}

func TestGetCities(t *testing.T) {
	for _, test := range getCitiesTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.Empty{}
		repositoryMock.On("GetCities", context.Background(), in).Return(&eventGrpc.GetCitiesRequest{
			Cities: test.outputRes,
		}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetCities()
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
