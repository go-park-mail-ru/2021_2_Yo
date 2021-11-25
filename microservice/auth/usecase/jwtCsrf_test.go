package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {

	useCaseTest := NewService(nil,nil)

	ctx := context.Background()
	protoUserId := &protoAuth.UserId{
		ID: "1",
	}
	_, err := useCaseTest.CreateToken(ctx,protoUserId)
	assert.NoError(t, err)
}

func TestCheckToken(t *testing.T) {
	useCaseTest := NewService(nil,nil)

	ctx := context.Background()
	userId := "1"
	csrfToken, err := generateCsrfToken(userId)
	assert.NoError(t,err)

	protoCSRF := &protoAuth.CSRFToken{
		CSRFToken: csrfToken,
	}
	protoUserId, err := useCaseTest.CheckToken(ctx,protoCSRF)

	assert.Equal(t,"1",protoUserId.ID)
	assert.NoError(t, err)
}
