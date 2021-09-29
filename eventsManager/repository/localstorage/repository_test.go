package localstorage

import (
	"backend/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestList(t *testing.T) {
	repositoryTest := NewRepositoryEventLocalStorage()
	resultEvents := make([]*models.Event, len(eventsDemo))
	for i := 0; i < len(eventsDemo); i++ {
		resultEvents[i] = toModelEvent(eventsDemo[i])
	}
	eventsTest, err := repositoryTest.List()
	require.NoError(t, err, "TestGetUser : repository.GetUser err = ", err)
	require.Equal(t, resultEvents, eventsTest, "TestList : lists of events are not equal")
}