package usecase

import (
	"backend/eventsManager/repository/localstorage"
	"backend/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestList(t *testing.T) {
	repositoryMock := new(localstorage.RepositoryEventMock)
	eventsManagerTest := NewUseCaseEvents(repositoryMock)
	var expected []*models.Event
	expected = nil
	repositoryMock.On("List")
	result, err := eventsManagerTest.List()
	require.NoError(t, err, "TestList : err = ", err)
	require.Equal(t, expected, result, "TestList : result and expected are not equal")
}
