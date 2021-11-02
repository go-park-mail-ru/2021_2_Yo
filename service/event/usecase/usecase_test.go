package usecase

import (
	"backend/service/event/repository/mock"
	"backend/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestList(t *testing.T) {
	repositoryMock := new(mock.RepositoryMock)
	eventsManagerTest := NewUseCase(repositoryMock)
	var expected []*models.Event
	expected = nil
	repositoryMock.On("List")
	result, err := eventsManagerTest.List()
	require.NoError(t, err, "TestList : err = ", err)
	require.Equal(t, expected, result, "TestList : result and expected are not equal")
}
