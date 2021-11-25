package session

import (
	"testing"
	"time"

	authServiceModels "backend/microservice/auth/models"

	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	client *redis.Client
)

var (
	key = "key"
	val = "val"
)

func TestCreate(t *testing.T) {
	exp := time.Duration(0)

	mock := redismock.NewNiceMock(client)

	mock.On("Set", key, val, exp).Return(redis.NewStatusResult("",nil))

	r := NewRepository(mock)

	data := &authServiceModels.SessionData{
		SessionId: key,
		UserId: val,
		Expiration: exp,
	}

	err := r.Create(data)
	assert.NoError(t, err)
}

func TestCheck(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	mock.On("Get", key).Return(redis.NewStringResult(val, nil))

	r := NewRepository(mock)
	res, err := r.Check(key)
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}
func TestDelete(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	keys := []string{key}
	mock.On("Del",keys).Return(redis.NewIntResult(0,nil))
	
	r := NewRepository(mock)

	err := r.Delete(key)
	assert.NoError(t, err)
}