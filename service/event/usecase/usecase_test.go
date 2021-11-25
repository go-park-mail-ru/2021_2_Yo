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
	tags      []string
	outputErr error
	outputRes []*models.Event
}{
	{1,
		"test",
		"test",
		[]string{"test"},
		nil,
		[]*models.Event{},
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &eventGrpc.GetEventsRequest{
			Title:    test.title,
			Category: test.category,
			Tags:     test.tags,
		}
		repositoryMock.On("GetEvents", context.Background(), in).Return(&eventGrpc.Events{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetEvents(test.title, test.category, test.tags)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes)
	}
}

/*

var updateEventTests = []struct {
	id            int
	event         *models.Event
	userId        string
	repositoryErr error
	outputErr     error
}{
	{1,
		&models.Event{
			ID: "1",
		},
		"1",
		nil,
		nil,
	},
	{2,
		&models.Event{},
		"1",
		nil,
		error2.ErrEmptyData,
	},
	{3,
		&models.Event{
			ID: "1",
		},
		"",
		nil,
		error2.ErrEmptyData,
	},
	{4,
		&models.Event{
			ID: "1",
		},
		"1",
		error2.ErrPostgres,
		error2.ErrPostgres,
	},
}

func TestUpdateEvent(t *testing.T) {
	for _, test := range updateEventTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("UpdateEvent", test.event, test.userId).Return(test.outputErr)
		actualErr := useCaseTest.UpdateEvent(test.event, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var deleteEventTests = []struct {
	id            int
	eventId       string
	userId        string
	repositoryErr error
	outputErr     error
}{
	{1,
		"1",
		"1",
		nil,
		nil,
	},
	{2,
		"",
		"1",
		nil,
		error2.ErrEmptyData,
	},
	{3,
		"1",
		"",
		nil,
		error2.ErrEmptyData,
	},
	{4,
		"1",
		"1",
		error2.ErrPostgres,
		error2.ErrPostgres,
	},
}

func TestDeleteEvent(t *testing.T) {
	for _, test := range deleteEventTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("DeleteEvent", test.eventId, test.userId).Return(test.outputErr)
		actualErr := useCaseTest.DeleteEvent(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getEventByIdTests = []struct {
	id            int
	eventId       string
	event         *models.Event
	repositoryErr error
	outputEvent   *models.Event
	outputErr     error
}{
	{1,
		"1",
		nil,
		nil,
		nil,
		nil,
	},
	{2,
		"",
		&models.Event{},
		nil,
		nil,
		error2.ErrEmptyData,
	},
	{3,
		"1",
		&models.Event{},
		error2.ErrPostgres,
		&models.Event{},
		error2.ErrPostgres,
	},
}

func TestGetEventById(t *testing.T) {
	for _, test := range getEventByIdTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetEventById", test.eventId).Return(test.event, test.repositoryErr)
		actualEvent, actualErr := useCaseTest.GetEventById(test.eventId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputEvent, actualEvent, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getEventsTests = []struct {
	id            int
	title         string
	category      string
	tags          []string
	repositoryErr error
	events        []*models.Event
	outputEvents  []*models.Event
	outputErr     error
}{
	{1,
		"testTitle",
		"testCategory",
		[]string{"testTag"},
		nil,
		nil,
		nil,
		nil,
	},
	{2,
		"testTitle",
		"testCategory",
		[]string{""},
		nil,
		nil,
		nil,
		nil,
	},
	{2,
		"testTitle",
		"testCategory",
		[]string{""},
		error2.ErrPostgres,
		nil,
		nil,
		error2.ErrPostgres,
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		tags := test.tags
		if test.tags != nil && test.tags[0] == "" {
			tags = nil
		}
		repositoryMock.On("GetEvents", test.title, test.category, tags).Return(test.events, test.repositoryErr)
		actualEvents, actualErr := useCaseTest.GetEvents(test.title, test.category, test.tags)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputEvents, actualEvents, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getEventsFromAuthorTests = []struct {
	id            int
	authorId      string
	repositoryErr error
	events        []*models.Event
	outputEvents  []*models.Event
	outputErr     error
}{
	{1,
		"1",
		nil,
		nil,
		nil,
		nil,
	},
	{2,
		"",
		nil,
		nil,
		nil,
		error2.ErrEmptyData,
	},
	{3,
		"1",
		error2.ErrPostgres,
		nil,
		nil,
		error2.ErrPostgres,
	},
}

func TestGetFromAuthorEvents(t *testing.T) {
	for _, test := range getEventsFromAuthorTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetEventsFromAuthor", test.authorId).Return(test.events, test.repositoryErr)
		actualEvents, actualErr := useCaseTest.GetEventsFromAuthor(test.authorId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputEvents, actualEvents, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

*/
