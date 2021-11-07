package postgres

import (
	"backend/models"
	error2 "backend/service/event/error"
	sql2 "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:event:repository:postgres:"

var createEventTests = []struct {
	id                    int
	event                 *models.Event
	expectationsShouldMet bool
	eventId               int
	postgresErr           error
	outputStr             string
	outputErr             error
}{
	{
		1,
		&models.Event{},
		true,
		100,
		nil,
		"100",
		nil,
	},
	{
		2,
		&models.Event{AuthorId: "incorrectAuthorId"},
		false,
		100,
		nil,
		"",
		error2.ErrAtoi,
	},
	{
		3,
		&models.Event{},
		true,
		100,
		sql2.ErrNoRows,
		"",
		error2.ErrNoRows,
	},
}

func TestCreateEvent(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range createEventTests {

		newEvent := &Event{}
		if test.event != nil {
			newEvent, err = toPostgresEvent(test.event)
			if err != error2.ErrAtoi {
				require.NoError(t, err)
			} else {
				newEvent = &Event{}
			}
		}

		//В этом случае не будут expectations
		if test.expectationsShouldMet {
			mock.ExpectQuery(createEventQuery).
				WithArgs(newEvent.Title,
					newEvent.Description,
					newEvent.Text,
					newEvent.City,
					newEvent.Category,
					newEvent.Viewed,
					newEvent.ImgUrl,
					newEvent.Date,
					newEvent.Geo,
					newEvent.Tag,
					newEvent.AuthorID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.eventId)).
				WillReturnError(test.postgresErr)
		}

		actualEventId, actualErr := repositoryTest.CreateEvent(test.event)
		t.Log(actualErr)
		t.Log(actualEventId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputStr, actualEventId, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateEventTests = []struct {
	id                    int
	event                 *models.Event
	userId                string
	expectationsShouldMet bool
	eventId               int
	postgresErr           error
	outputErr             error
}{
	{
		1,
		&models.Event{ID: "1", AuthorId: "100"},
		"100",
		true,
		100,
		nil,
		nil,
	},
}

func TestUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range updateEventTests {

		newEvent := &Event{}
		if test.event != nil {
			newEvent, err = toPostgresEvent(test.event)
			if err != error2.ErrAtoi {
				require.NoError(t, err)
			} else {
				newEvent = &Event{}
			}
		}

		//В этом случае не будут expectations
		if test.expectationsShouldMet {
			mock.ExpectQuery(checkAuthorQuery).
				WithArgs(newEvent.ID).
				WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(test.event.AuthorId)).
				WillReturnError(test.postgresErr)
			mock.ExpectQuery(updateEventQuery).
				WithArgs(newEvent.Title,
					newEvent.Description,
					newEvent.Text,
					newEvent.City,
					newEvent.Category,
					newEvent.Viewed,
					newEvent.ImgUrl,
					newEvent.Date,
					newEvent.Geo,
					newEvent.Tag,
					newEvent.ID).
				WillReturnError(test.postgresErr)
		}

		actualErr := repositoryTest.UpdateEvent(test.event, test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
