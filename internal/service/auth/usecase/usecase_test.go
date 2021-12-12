package usecase

import (
	protoAuth "backend/internal/microservice/auth/proto"
	"backend/internal/microservice/auth/usecase"
	"backend/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

var signUpTests = []struct {
	id        int
	input     *models.User
	clientRes *protoAuth.UserId
	clientErr error
	output    string
}{
	{
		1,
		&models.User{},
		&protoAuth.UserId{
			ID: "test",
		},
		nil,
		"test",
	},
	{
		2,
		&models.User{},
		&protoAuth.UserId{
			ID: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestSignUp(t *testing.T) {
	for _, test := range signUpTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.SignUpRequest{
			Name:     test.input.Name,
			Surname:  test.input.Surname,
			Mail:     test.input.Mail,
			Password: test.input.Password,
		}
		clientMock.On("SignUp", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.SignUp(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}

var signInTests = []struct {
	id        int
	input     *models.User
	clientRes *protoAuth.UserId
	clientErr error
	output    string
}{
	{
		1,
		&models.User{},
		&protoAuth.UserId{
			ID: "test",
		},
		nil,
		"test",
	},
	{
		2,
		&models.User{},
		&protoAuth.UserId{
			ID: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestSignIn(t *testing.T) {
	for _, test := range signUpTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.SignInRequest{
			Mail:     test.input.Mail,
			Password: test.input.Password,
		}
		clientMock.On("SignIn", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.SignIn(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}

var createSessionTests = []struct {
	id        int
	input     string
	clientRes *protoAuth.Session
	clientErr error
	output    string
}{
	{
		1,
		"test",
		&protoAuth.Session{
			Session: "test",
		},
		nil,
		"test",
	},
	{
		2,
		"test",
		&protoAuth.Session{
			Session: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestCreateSession(t *testing.T) {
	for _, test := range createSessionTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.UserId{
			ID: test.input,
		}
		clientMock.On("CreateSession", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.CreateSession(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}

var checkSessionTests = []struct {
	id        int
	input     string
	clientRes *protoAuth.UserId
	clientErr error
	output    string
}{
	{
		1,
		"test",
		&protoAuth.UserId{
			ID: "test",
		},
		nil,
		"test",
	},
	{
		2,
		"test",
		&protoAuth.UserId{
			ID: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestCheckSession(t *testing.T) {
	for _, test := range checkSessionTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.Session{
			Session: test.input,
		}
		clientMock.On("CheckSession", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.CheckSession(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}

var deleteSessionTests = []struct {
	id        int
	input     string
	clientRes *protoAuth.UserId
	clientErr error
}{
	{
		1,
		"test",
		&protoAuth.UserId{
			ID: "test",
		},
		nil,
	},
	{
		2,
		"test",
		&protoAuth.UserId{
			ID: "",
		},
		errors.New("test_err"),
	},
}

func TestDeleteSession(t *testing.T) {
	for _, test := range deleteSessionTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.Session{
			Session: test.input,
		}
		clientMock.On("DeleteSession", context.Background(), in).Return(&protoAuth.Success{}, test.clientErr)
		err := useCaseTest.DeleteSession(test.input)
		require.Equal(t, test.clientErr, err)
	}
}

var createTokenTests = []struct {
	id        int
	input     string
	clientRes *protoAuth.CSRFToken
	clientErr error
	output    string
}{
	{
		1,
		"test",
		&protoAuth.CSRFToken{
			CSRFToken: "test",
		},
		nil,
		"test",
	},
	{
		2,
		"test",
		&protoAuth.CSRFToken{
			CSRFToken: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestCreateToken(t *testing.T) {
	for _, test := range createTokenTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.UserId{
			ID: test.input,
		}
		clientMock.On("CreateToken", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.CreateToken(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}

var checkTokenTests = []struct {
	id        int
	input     string
	clientRes *protoAuth.UserId
	clientErr error
	output    string
}{
	{
		1,
		"test",
		&protoAuth.UserId{
			ID: "test",
		},
		nil,
		"test",
	},
	{
		2,
		"test",
		&protoAuth.UserId{
			ID: "",
		},
		errors.New("test_err"),
		"",
	},
}

func TestCheckToken(t *testing.T) {
	for _, test := range checkSessionTests {
		clientMock := new(usecase.AuthClientMock)
		useCaseTest := NewUseCase(clientMock)
		in := &protoAuth.CSRFToken{
			CSRFToken: test.input,
		}
		clientMock.On("CheckToken", context.Background(), in).Return(test.clientRes, test.clientErr)
		res, err := useCaseTest.CheckToken(test.input)
		require.Equal(t, test.clientErr, err)
		require.Equal(t, test.output, res)
	}
}
