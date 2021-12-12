package usecase

import (
	"backend/internal/microservice/user/proto"
	error2 "backend/internal/service/user/error"
	"backend/pkg/models"
	"backend/pkg/utils"
	"context"
	"errors"

	//"backend/pkg/utils"
	"backend/microservice/user/repository"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:user:usecase:"

var getUserByIdTests = []struct {
	id         int
	input      string
	outputUser *models.User
	outputErr  error
}{
	{
		1,
		"1",
		&models.User{},
		nil,
	},
	{
		2,
		"",
		nil,
		error2.ErrEmptyData,
	},
	{
		1,
		"1",
		nil,
		errors.New("test_err"),
	},
}

func TestGetUserById(t *testing.T) {
	for _, test := range getUserByIdTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.UserId{
			ID: test.input,
		}
		repositoryMock.On("GetUserById", context.Background(), in).Return(&userGrpc.User{}, test.outputErr)
		actualUser, actualErr := useCaseTest.GetUserById(test.input)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputUser, actualUser, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateUserInfoTests = []struct {
	id        int
	input     *models.User
	outputErr error
}{
	{
		1,
		&models.User{
			ID:      "test",
			Name:    "test",
			Surname: "test",
		},
		nil,
	},
	{
		2,
		&models.User{},
		error2.ErrEmptyData,
	},
	{
		3,
		&models.User{
			ID:      "test",
			Name:    "test",
			Surname: "test",
		},
		errors.New("test_err"),
	},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := MakeProtoUser(test.input)
		repositoryMock.On("UpdateUserInfo", context.Background(), in).Return(&userGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.UpdateUserInfo(test.input)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateUserPasswordTests = []struct {
	id        int
	userId    string
	password  string
	outputErr error
}{
	{
		1,
		"test",
		"test",
		nil,
	},
	{
		2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{
		3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestUpdateUserPassword(t *testing.T) {
	for _, test := range updateUserPasswordTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.UpdateUserPasswordRequest{
			ID:       test.userId,
			Password: utils.CreatePasswordHash(test.password),
		}
		repositoryMock.On("UpdateUserPassword", context.Background(), in).Return(&userGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.UpdateUserPassword(test.userId, test.password)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getSubscribersTests = []struct {
	id        int
	userId    string
	outputErr error
	outputRes []*models.User
}{
	{
		1,
		"test",
		nil,
		[]*models.User{},
	},
	{
		2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{
		3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetSubscribers(t *testing.T) {
	for _, test := range getSubscribersTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.UserId{
			ID: test.userId,
		}
		repositoryMock.On("GetSubscribers", context.Background(), in).Return(&userGrpc.Users{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetSubscribers(test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getSubscribesTests = []struct {
	id        int
	userId    string
	outputErr error
	outputRes []*models.User
}{
	{
		1,
		"test",
		nil,
		[]*models.User{},
	},
	{
		2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{
		3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetSubscribes(t *testing.T) {
	for _, test := range getSubscribesTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.UserId{
			ID: test.userId,
		}
		repositoryMock.On("GetSubscribes", context.Background(), in).Return(&userGrpc.Users{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetSubscribes(test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getVisitorsTests = []struct {
	id        int
	eventId   string
	outputErr error
	outputRes []*models.User
}{
	{
		1,
		"test",
		nil,
		[]*models.User{},
	},
	{
		2,
		"",
		error2.ErrEmptyData,
		nil,
	},
	{
		3,
		"test",
		errors.New("test_err"),
		nil,
	},
}

func TestGetVisitors(t *testing.T) {
	for _, test := range getVisitorsTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.EventId{
			ID: test.eventId,
		}
		repositoryMock.On("GetVisitors", context.Background(), in).Return(&userGrpc.Users{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetVisitors(test.eventId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var subscribeTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	outputErr    error
}{
	{
		1,
		"test",
		"test",
		nil,
	},
	{
		2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{
		3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestSubscribe(t *testing.T) {
	for _, test := range subscribeTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		repositoryMock.On("Subscribe", context.Background(), in).Return(&userGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.Subscribe(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var unsubscribeTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	outputErr    error
}{
	{
		1,
		"test",
		"test",
		nil,
	},
	{
		2,
		"",
		"",
		error2.ErrEmptyData,
	},
	{
		3,
		"test",
		"test",
		errors.New("test_err"),
	},
}

func TestUnsubscribe(t *testing.T) {
	for _, test := range unsubscribeTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		repositoryMock.On("Unsubscribe", context.Background(), in).Return(&userGrpc.Empty{}, test.outputErr)
		actualErr := useCaseTest.Unsubscribe(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var isSubscribedTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	outputErr    error
	outputRes    bool
}{
	{
		1,
		"test1",
		"test2",
		nil,
		true,
	},
	{
		2,
		"",
		"",
		error2.ErrEmptyData,
		false,
	},
	{
		3,
		"test",
		"test",
		nil,
		false,
	},
	{
		4,
		"test1",
		"test2",
		errors.New("test_err"),
		false,
	},
}

func TestIsSubscribed(t *testing.T) {
	for _, test := range isSubscribedTests {
		repositoryMock := new(repository.RepositoryClientMock)
		useCaseTest := NewUseCase(repositoryMock)
		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		repositoryMock.On("IsSubscribed", context.Background(), in).Return(&userGrpc.IsSubscribedRequest{
			Result: test.outputRes,
		}, test.outputErr)
		actualRes, actualErr := useCaseTest.IsSubscribed(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputRes, actualRes, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

/*

var updateUserInfoTests = []struct {
	id        int
	user      *models.User
	outputErr error
}{
	{
		1,
		&models.User{
			ID:      "1",
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		nil,
	},
	{
		2,
		&models.User{
			ID:      "",
			Name:    "",
			Surname: "",
			About:   "",
		},
		error2.ErrEmptyData,
	},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("UpdateUserInfo", test.user).Return(test.outputErr)
		actualErr := useCaseTest.UpdateUserInfo(test.user)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateUserPasswordTests = []struct {
	id        int
	userId    string
	password  string
	outputErr error
}{
	{
		1,
		"1",
		"testPassword",
		nil,
	},
	{
		2,
		"",
		"",
		error2.ErrEmptyData,
	},
}

func TestUpdateUserPassword(t *testing.T) {
	for _, test := range updateUserPasswordTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		hashedPassword := utils.CreatePasswordHash(test.password)
		repositoryMock.On("UpdateUserPassword", test.userId, hashedPassword).Return(test.outputErr)
		actualErr := useCaseTest.UpdateUserPassword(test.userId, test.password)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

*/
