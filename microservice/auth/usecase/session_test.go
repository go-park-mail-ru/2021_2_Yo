package usecase

import (
	authServiceModels "backend/microservice/auth/models"
	protoAuth "backend/microservice/auth/proto"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	sessionRepositoryMock := new(AuthSessionMock)
	authRepositoryMock := new(AuthRepoMock)
	useCaseTest := NewService(authRepositoryMock, sessionRepositoryMock)
	userId := "-1"
	sessionData := &authServiceModels.SessionData{
		SessionId:  "",
		UserId:     userId,
		Expiration: sessionLifeTime,
	}
	sessionRepositoryMock.On("Create", sessionData).Return(nil)
	ctx := context.Background()
	protoUserId := &protoAuth.UserId{
		ID: userId,
	}
	protoSession, err := useCaseTest.CreateSession(ctx, protoUserId)

	assert.Equal(t, len(protoSession.Session), 0)
	assert.NoError(t, err)
	sessionRepositoryMock.AssertExpectations(t)
}

func TestCheckSession(t *testing.T) {
	sessionRepositoryMock := new(AuthSessionMock)

	sessionId := "1111111111111111"
	expUserId := "1"
	sessionRepositoryMock.On("Check", sessionId).Return(expUserId, nil)

	useCaseTest := NewService(nil, sessionRepositoryMock)

	ctx := context.Background()
	protoSession := &protoAuth.Session{
		Session: sessionId,
	}
	protoUserId, err := useCaseTest.CheckSession(ctx, protoSession)
	userId := protoUserId.ID
	assert.Equal(t, "1", userId)
	assert.NoError(t, err)
	sessionRepositoryMock.AssertExpectations(t)
}

func TestDeleteSession(t *testing.T) {
	sessionRepositoryMock := new(AuthSessionMock)

	sessionId := "1111111111111111"

	sessionRepositoryMock.On("Delete", sessionId).Return(nil)

	useCaseTest := NewService(nil, sessionRepositoryMock)

	ctx := context.Background()
	protoSession := &protoAuth.Session{
		Session: sessionId,
	}
	protoSuccess, err := useCaseTest.DeleteSession(ctx, protoSession)
	success := protoSuccess.Ok
	assert.Equal(t, "success", success)
	assert.NoError(t, err)
	sessionRepositoryMock.AssertExpectations(t)
}
