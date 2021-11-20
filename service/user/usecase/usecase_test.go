package usecase

import (
	"backend/pkg/models"
	"backend/pkg/utils"
	error2 "backend/service/user/error"
	"backend/service/user/repository/mock"
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
		&models.User{
			ID: "1",
		},
		nil,
	},
	{
		2,
		"",
		nil,
		error2.ErrEmptyData,
	},
}

func TestGetUserById(t *testing.T) {
	for _, test := range getUserByIdTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock)
		repositoryMock.On("GetUserById", test.input).Return(test.outputUser, test.outputErr)
		actualUser, actualErr := useCaseTest.GetUserById(test.input)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputUser, actualUser, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

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
