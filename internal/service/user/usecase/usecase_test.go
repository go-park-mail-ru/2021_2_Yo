package usecase

import (
	"backend/internal/models"
	error2 "backend/internal/service/user/error"
	"backend/internal/service/user/repository/mock"
	"backend/internal/utils"
	"errors"

	"github.com/stretchr/testify/require"
	"testing"
)

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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetUserById", test.input).Return(&models.User{}, test.outputErr)
		actualUser, actualErr := useCaseTest.GetUserById(test.input)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputUser, actualUser)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("UpdateUserInfo", test.input).Return(test.outputErr)
		actualErr := useCaseTest.UpdateUserInfo(test.input)
		require.Equal(t, test.outputErr, actualErr)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("UpdateUserPassword", test.userId, utils.CreatePasswordHash(test.password)).Return(test.outputErr)
		actualErr := useCaseTest.UpdateUserPassword(test.userId, test.password)
		require.Equal(t, test.outputErr, actualErr)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetSubscribers", test.userId).Return([]*models.User{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetSubscribers(test.userId)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputRes, actualRes)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetSubscribes", test.userId).Return([]*models.User{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetSubscribes(test.userId)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputRes, actualRes)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetVisitors", test.eventId).Return([]*models.User{}, test.outputErr)
		actualRes, actualErr := useCaseTest.GetVisitors(test.eventId)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputRes, actualRes)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("Subscribe", test.subscribedId, test.subscriberId).Return(test.outputErr)
		actualErr := useCaseTest.Subscribe(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("Unsubscribe", test.subscribedId, test.subscriberId).Return(test.outputErr)
		actualErr := useCaseTest.Unsubscribe(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr)
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
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("IsSubscribed", test.subscribedId, test.subscriberId).Return(test.outputRes, test.outputErr)
		actualRes, actualErr := useCaseTest.IsSubscribed(test.subscribedId, test.subscriberId)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputRes, actualRes)
	}
}
