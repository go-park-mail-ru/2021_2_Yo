package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {

	sessionId := "1111111111111111"
	expUserId := "1"

	useCaseTest := NewService(nil,nil)

	ctx := context.Background()
	protoSession := &protoAuth.CSRFToken{
		Session: sessionId,
	}
	protoUserId, err := useCaseTest.CreateToken(ctx,protoSession)
	userId := protoUserId.ID
	assert.Equal(t, "1",userId)
	assert.NoError(t, err)
	sessionRepositoryMock.AssertExpectations(t)
}
