package eventRepository

import (
	eventGrpc "backend/microservice/event/proto"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	error3 "backend/service/user/error"
	"context"
	sql2 "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

var createEventTests = []struct {
	id          int
	event       *models.Event
	eventId     int
	postgresErr error
	output      string
	outputErr   error
}{
	{
		1,
		&models.Event{
			AuthorId: "1",
		},
		10,
		nil,
		"10",
		nil,
	},
	{
		2,
		&models.Event{},
		10,
		nil,
		"10",
		nil,
	},
	{
		3,
		&models.Event{
			AuthorId: "test",
		},
		0,
		nil,
		"",
		error2.ErrAtoi,
	},
	{
		4,
		&models.Event{
			AuthorId: "10",
		},
		0,
		sql2.ErrNoRows,
		"",
		error2.ErrNoRows,
	},
	{
		5,
		&models.Event{
			AuthorId: "10",
		},
		0,
		sql2.ErrConnDone,
		"",
		error2.ErrPostgres,
	},
}

func TestCreateEvent(t *testing.T) {
	for _, test := range createEventTests {

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		newEvent, err := toPostgresEvent(test.event)
		if err != nil {
			newEvent = &Event{}
		}

		mock.ExpectQuery(createEventQuery).
			WithArgs(
				newEvent.Title,
				newEvent.Description,
				newEvent.Text,
				newEvent.City,
				newEvent.Category,
				newEvent.Viewed,
				newEvent.ImgUrl,
				newEvent.Date,
				newEvent.Geo,
				newEvent.Address,
				newEvent.Tag,
				newEvent.AuthorID,
			).WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(test.eventId)).WillReturnError(test.postgresErr)

		in := toProtoEvent(test.event)
		out, actualErr := repositoryTest.CreateEvent(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := out.ID
		require.Equal(t, test.output, actualRes)
	}
}

var updateEventTests = []struct {
	id          int
	event       *models.Event
	userId      string
	postgresErr error
	outputErr   error
}{
	{
		1,
		&models.Event{
			ID:       "10",
			AuthorId: "10",
		},
		"10",
		nil,
		nil,
	},
	{
		2,
		&models.Event{
			ID:       "1",
			AuthorId: "10",
			ImgUrl:   "test",
		},
		"10",
		nil,
		nil,
	},
	{
		3,
		&models.Event{
			ID:       "1",
			AuthorId: "10",
			ImgUrl:   "test",
		},
		"10",
		error2.ErrPostgres,
		error2.ErrPostgres,
	},
	{
		4,
		&models.Event{
			ID:       "1",
			AuthorId: "10",
		},
		"10",
		error2.ErrPostgres,
		error2.ErrPostgres,
	},
	{
		5,
		&models.Event{
			ID:       "1",
			AuthorId: "10",
		},
		"2",
		nil,
		error2.ErrNotAllowed,
	},
	{
		6,
		&models.Event{
			ID:       "test",
			AuthorId: "10",
		},
		"2",
		nil,
		error2.ErrAtoi,
	},
	{
		7,
		&models.Event{
			ID:       "10",
			AuthorId: "10",
		},
		"test",
		nil,
		error2.ErrAtoi,
	},
}

func TestUpdateEvent(t *testing.T) {
	for _, test := range updateEventTests {

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		newEvent, err := toPostgresEvent(test.event)
		if err != nil {
			newEvent = &Event{}
		}

		eventIdInt, err := strconv.Atoi(test.event.ID)
		if err != nil {
			eventIdInt = 0
		}
		newEvent.ID = eventIdInt

		authorIdInt, err := strconv.Atoi(test.event.AuthorId)
		if err != nil {
			eventIdInt = 0
		}

		mock.ExpectQuery(checkAuthorQuery).
			WithArgs(eventIdInt).
			WillReturnRows(sqlmock.NewRows([]string{"author_id"}).
				AddRow(authorIdInt)).
			WillReturnError(nil)

		if test.event.ImgUrl != "" {
			mock.ExpectQuery(updateEventQuery).
				WithArgs(
					newEvent.Title,
					newEvent.Description,
					newEvent.Text,
					newEvent.City,
					newEvent.Category,
					newEvent.Viewed,
					newEvent.ImgUrl,
					newEvent.Date,
					newEvent.Geo,
					newEvent.Address,
					newEvent.Tag,
					newEvent.ID,
				).WillReturnRows(sqlmock.NewRows([]string{""})).WillReturnError(test.postgresErr)
		} else {
			mock.ExpectQuery(updateEventQueryWithoutImgUrl).
				WithArgs(
					newEvent.Title,
					newEvent.Description,
					newEvent.Text,
					newEvent.City,
					newEvent.Category,
					newEvent.Viewed,
					newEvent.Date,
					newEvent.Geo,
					newEvent.Address,
					newEvent.Tag,
					newEvent.ID,
				).WillReturnRows(sqlmock.NewRows([]string{""})).WillReturnError(test.postgresErr)
		}

		protoEvent := toProtoEvent(test.event)
		in := &eventGrpc.UpdateEventRequest{
			Event:  protoEvent,
			UserId: test.userId,
		}
		_, actualErr := repositoryTest.UpdateEvent(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var deleteEventTests = []struct {
	id          int
	eventId     string
	userId      string
	authorId    int
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"1",
		1,
		nil,
		nil,
	},
	{
		2,
		"test",
		"1",
		1,
		nil,
		error2.ErrAtoi,
	},
	{
		3,
		"1",
		"test",
		1,
		nil,
		error2.ErrAtoi,
	},
	{
		4,
		"1",
		"1",
		10,
		nil,
		error2.ErrNotAllowed,
	},
	{
		5,
		"1",
		"1",
		1,
		error2.ErrPostgres,
		error2.ErrPostgres,
	},
}

func TestDeleteEvent(t *testing.T) {

	for _, test := range deleteEventTests {

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		eventIdInt, err := strconv.Atoi(test.eventId)
		if err != nil {
			eventIdInt = 0
		}

		mock.ExpectQuery(checkAuthorQuery).
			WithArgs(eventIdInt).
			WillReturnRows(sqlmock.NewRows([]string{"author_id"}).
				AddRow(test.authorId)).
			WillReturnError(nil)

		mock.ExpectQuery(deleteEventQuery).
			WithArgs(eventIdInt).WillReturnRows(sqlmock.NewRows([]string{})).WillReturnError(test.postgresErr)

		in := &eventGrpc.DeleteEventRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		_, actualErr := repositoryTest.DeleteEvent(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var getEventByIdTests = []struct {
	id          int
	eventId     string
	postgresErr error
	outputEvent *models.Event
	outputErr   error
}{
	{
		1,
		"1",
		nil,
		&models.Event{
			ID:       "0",
			Title:    "test",
			AuthorId: "0",
		},
		nil,
	},
	{
		1,
		"1",
		sql2.ErrNoRows,
		&models.Event{},
		error2.ErrNoRows,
	},
	{
		1,
		"1",
		sql2.ErrConnDone,
		&models.Event{},
		error2.ErrPostgres,
	},
}

func TestGetEventById(t *testing.T) {
	for _, test := range getEventByIdTests {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		mock.ExpectQuery(getEventQuery).WithArgs(test.eventId).
			WillReturnRows(sqlmock.NewRows([]string{"title"}).
				AddRow(test.outputEvent.Title)).WillReturnError(test.postgresErr)

		in := &eventGrpc.EventId{ID: test.eventId}
		out, actualErr := repositoryTest.GetEventById(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualEvent := fromProtoToModel(out)
		require.Equal(t, test.outputEvent, actualEvent)
	}
}

var getEventsTests = []struct {
	id           int
	title        string
	category     string
	tags         []string
	postgresErr  error
	outputEvents []*models.Event
	outputErr    error
}{
	{
		id:          1,
		title:       "test",
		category:    "test",
		tags:        []string{"test"},
		postgresErr: nil,
		outputEvents: []*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
		outputErr: nil,
	},
	{
		id:          2,
		title:       "",
		category:    "",
		tags:        nil,
		postgresErr: nil,
		outputEvents: []*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
		outputErr: nil,
	},
	{
		id:           3,
		title:        "",
		category:     "",
		tags:         nil,
		postgresErr:  sql2.ErrNoRows,
		outputEvents: []*models.Event{},
		outputErr:    sql2.ErrNoRows,
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		postgresTags := make(pq.StringArray, len(test.tags))
		for i := range test.tags {
			postgresTags[i] = test.tags[i]
		}
		query := listQuery + " "
		if test.title != "" {
			query += `where lower(title) ~ lower($1) and `
		} else {
			query += `where $1 = $1 and `
		}
		if test.category != "" {
			query += `lower(category) = lower($2) and `
		} else {
			query += `$2 = $2 and `
		}
		if len(postgresTags) != 0 {
			query += `tag && $3::varchar[]`
		} else {
			query += `$3 = $3`
		}
		query += "order by viewed DESC"

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(query).
			WithArgs(test.title, test.category, postgresTags).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.GetEventsRequest{
			Title:    test.title,
			Category: test.category,
			Tags:     test.tags,
		}
		out, actualErr := repositoryTest.GetEvents(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.Event, len(out.Events))
		for i, e := range out.Events {
			actualRes[i] = fromProtoToModel(e)
		}
		require.Equal(t, test.outputEvents, actualRes)
	}
}

var getVisitedEventsTests = []struct {
	id          int
	userId      string
	postgresErr error
	outputErr   error
	outputRes   []*models.Event
}{
	{
		1,
		"1",
		nil,
		nil,
		[]*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
	},
	{
		2,
		"a",
		nil,
		error3.ErrAtoi,
		[]*models.Event{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error3.ErrPostgres,
		[]*models.Event{},
	},
}

func TestGetVisitedEvents(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getVisitedEventsTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(visitedQuery).
			WithArgs(userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.UserId{
			ID: test.userId,
		}
		out, actualErr := repositoryTest.GetVisitedEvents(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.Event, len(out.Events))
		for i, e := range out.Events {
			actualRes[i] = fromProtoToModel(e)
		}
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getCreatedEventsTests = []struct {
	id          int
	userId      string
	postgresErr error
	outputErr   error
	outputRes   []*models.Event
}{
	{
		1,
		"1",
		nil,
		nil,
		[]*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
	},
	{
		2,
		"a",
		nil,
		error3.ErrAtoi,
		[]*models.Event{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error3.ErrPostgres,
		[]*models.Event{},
	},
}

func TestGetCreatedEvents(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getCreatedEventsTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(createdQuery).
			WithArgs(userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.UserId{
			ID: test.userId,
		}
		out, actualErr := repositoryTest.GetCreatedEvents(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.Event, len(out.Events))
		for i, e := range out.Events {
			actualRes[i] = fromProtoToModel(e)
		}
		require.Equal(t, test.outputRes, actualRes)
	}
}

var visitTests = []struct {
	id          int
	eventId     string
	userId      string
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"2",
		nil,
		nil,
	},
	{
		2,
		"a",
		"b",
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		nil,
		error3.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestVisit(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range visitTests {
		eventIdInt, err := strconv.Atoi(test.eventId)
		if err != nil {
			eventIdInt = 0
		}
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(visitQuery).
			WithArgs(eventIdInt, userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		_, actualErr := repositoryTest.Visit(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var unvisitTests = []struct {
	id          int
	eventId     string
	userId      string
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"2",
		nil,
		nil,
	},
	{
		2,
		"a",
		"b",
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		nil,
		error3.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestUnvisit(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range unvisitTests {
		eventIdInt, err := strconv.Atoi(test.eventId)
		if err != nil {
			eventIdInt = 0
		}
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(unvisitQuery).
			WithArgs(eventIdInt, userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		_, actualErr := repositoryTest.Unvisit(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var isVisitedTests = []struct {
	id          int
	eventId     string
	userId      string
	count       int
	result      bool
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"2",
		10,
		true,
		nil,
		nil,
	},
	{
		1,
		"1",
		"2",
		0,
		false,
		nil,
		nil,
	},
	{
		2,
		"a",
		"b",
		0,
		false,
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		0,
		false,
		nil,
		error3.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		0,
		false,
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestIsVisited(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range isVisitedTests {
		eventIdInt, err := strconv.Atoi(test.eventId)
		if err != nil {
			eventIdInt = 0
		}
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(test.count)
		mock.ExpectQuery(isVisitedQuery).
			WithArgs(eventIdInt, userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &eventGrpc.VisitRequest{
			EventId: test.eventId,
			UserId:  test.userId,
		}
		out, actualErr := repositoryTest.IsVisited(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := out.Result
		require.Equal(t, test.result, actualRes)
	}
}
