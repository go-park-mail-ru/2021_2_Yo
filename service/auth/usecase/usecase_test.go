package usecase

import (
	"backend/models"
	error2 "backend/service/auth/error"
	"backend/service/auth/repository/mock"
	"backend/utils"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:auth:usecase:"

var signUpTests = []struct {
	id        int
	input     *models.User
	outputStr string
	outputErr error
}{
	{
		1,
		&models.User{
			ID: "",
		},
		"",
		nil,
	},
	{
		2,
		nil,
		"",
		error2.ErrEmptyData,
	},
}

func TestSignUp(t *testing.T) {
	secretWord := []byte("secret")
	for _, test := range signUpTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock, secretWord)
		repositoryMock.On("CreateUser", test.input).Return(test.outputStr, test.outputErr)
		actualStr, actualErr := useCaseTest.SignUp(test.input)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputStr, actualStr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var signInTests = []struct {
	id        int
	inputMail string
	inputPass string
	user      *models.User
	outputStr string
	outputErr error
}{
	{
		1,
		"a",
		"b",
		&models.User{
			ID: "1",
		},
		"1",
		nil,
	},
	{
		2,
		"",
		"",
		&models.User{
			ID: "1",
		},
		"",
		error2.ErrEmptyData,
	},
	{
		3,
		"a",
		"b",
		&models.User{
			ID: "1",
		},
		"",
		error2.ErrPostgres,
	},
}

func TestSignIn(t *testing.T) {
	secretWord := []byte("secret")
	for _, test := range signInTests {
		repositoryMock := new(mock.RepositoryMock)
		useCaseTest := NewUseCase(repositoryMock, secretWord)
		hashedPass := utils.CreatePasswordHash(test.inputPass)
		repositoryMock.On("GetUser", test.inputMail, hashedPass).Return(test.user, test.outputErr)
		actualStr, actualErr := useCaseTest.SignIn(test.inputMail, test.inputPass)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputStr, actualStr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
