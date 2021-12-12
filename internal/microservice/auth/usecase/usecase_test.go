package usecase

import (
	protoAuth "backend/internal/microservice/auth/proto"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignUp(t *testing.T) {
	authRepositoryMock := new(AuthRepoMock)

	newUser := &models.User{
		Name:     "Artyom",
		Surname:  "Shirshov",
		Mail:     "test@mail.ru",
		Password: utils.CreatePasswordHash("12345678"),
	}
	expUserId := "1"
	authRepositoryMock.On("CreateUser", newUser).Return(expUserId, nil)

	useCaseTest := NewService(authRepositoryMock, nil)

	ctx := context.Background()
	protoSignUp := &protoAuth.SignUpRequest{
		Name:     "Artyom",
		Surname:  "Shirshov",
		Mail:     "test@mail.ru",
		Password: "12345678",
	}
	protoUserId, err := useCaseTest.SignUp(ctx, protoSignUp)
	userId := protoUserId.ID
	assert.Equal(t, "1", userId)
	assert.NoError(t, err)
	authRepositoryMock.AssertExpectations(t)
}

func TestSignIn(t *testing.T) {
	authRepositoryMock := new(AuthRepoMock)
	password := "12345678"
	newUser := &models.User{
		ID:       "1",
		Name:     "Artyom",
		Surname:  "Shirshov",
		Mail:     "test@mail.ru",
		Password: utils.CreatePasswordHash(password),
	}

	authRepositoryMock.On("GetUser", newUser.Mail, newUser.Password).Return(newUser, nil)

	useCaseTest := NewService(authRepositoryMock, nil)

	ctx := context.Background()
	protoSignIn := &protoAuth.SignInRequest{
		Mail:     "test@mail.ru",
		Password: "12345678",
	}
	protoUserId, err := useCaseTest.SignIn(ctx, protoSignIn)
	userId := protoUserId.ID
	assert.Equal(t, "1", userId)
	assert.NoError(t, err)
	authRepositoryMock.AssertExpectations(t)
}
