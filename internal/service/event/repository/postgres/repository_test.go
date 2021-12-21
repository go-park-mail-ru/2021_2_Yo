package postgres

import (
	"backend/internal/models"
	error2 "backend/internal/service/event/error"
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
		out, actualErr := repositoryTest.CreateEvent(test.event)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.output, out)
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
					newEvent.Date,
					newEvent.Geo,
					newEvent.Address,
					newEvent.Tag,
					newEvent.ID,
				).WillReturnRows(sqlmock.NewRows([]string{""})).WillReturnError(test.postgresErr)
		}
		actualErr := repositoryTest.UpdateEvent(test.event, test.userId)
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

		actualErr := repositoryTest.DeleteEvent(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var getEventByIdTests = []struct {
	id          int
	eventId     string
	postgresErr error
	outputRes   *models.Event
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

		eventIdInt, _ := strconv.Atoi(test.eventId)

		mock.ExpectQuery(incrementEventViews).WithArgs(eventIdInt).
			WillReturnRows(sqlmock.NewRows([]string{"title"}).
				AddRow(test.outputRes.Title)).WillReturnError(nil)

		mock.ExpectQuery(getEventQuery).WithArgs(eventIdInt).
			WillReturnRows(sqlmock.NewRows([]string{"title"}).
				AddRow(test.outputRes.Title)).WillReturnError(test.postgresErr)

		out, actualErr := repositoryTest.GetEventById(test.eventId)
		require.Equal(t, test.outputErr, actualErr)
		if test.outputErr != nil {
			require.Equal(t, (*models.Event)(nil), out)
		} else {
			require.Equal(t, test.outputRes, out)
		}
	}
}

var getEventsTests = []struct {
	id          int
	userId      string
	title       string
	category    string
	city        string
	date        string
	tags        []string
	postgresErr error
	outputRes   []*models.Event
	outputErr   error
}{
	{
		id:          1,
		userId:      "",
		title:       "test",
		category:    "test",
		city:        "test",
		date:        "test",
		tags:        []string{"test"},
		postgresErr: nil,
		outputRes: []*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
		outputErr: nil,
	},
	{
		id:          2,
		userId:      "",
		title:       "",
		category:    "",
		city:        "",
		date:        "",
		tags:        nil,
		postgresErr: nil,
		outputRes: []*models.Event{
			&models.Event{
				ID:       "1",
				AuthorId: "0",
			},
		},
		outputErr: nil,
	},
	{
		id:          3,
		userId:      "",
		title:       "",
		category:    "",
		city:        "",
		date:        "",
		tags:        nil,
		postgresErr: sql2.ErrNoRows,
		outputRes:   []*models.Event{},
		outputErr:   error2.ErrPostgres,
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)

		var userIdInt int
		if test.userId == "" {
			userIdInt = 0
		} else {
			userIdInt1, _ := strconv.Atoi(test.userId)
			userIdInt = userIdInt1
		}

		postgresTags := make(pq.StringArray, len(test.tags))
		for i := range test.tags {
			postgresTags[i] = test.tags[i]
		}
		query := `select e.*, count(v) from event as e
				left join visitor as v on e.id = v.event_id and `
		query += `v.user_id = $1 `
		if test.title != "" {
			query += `where lower(title) ~ lower($2) and `
		} else {
			query += `where $2 = $2 and `
		}
		if test.category != "" {
			query += `lower(category) = lower($3) and `
		} else {
			query += `$3 = $3 and `
		}
		if test.city != "" {
			query += `lower(city) = lower($4) and `
		} else {
			query += `$4 = $4 and `
		}
		if test.date != "" {
			query += `lower(e.date) = lower($5) and `
		} else {
			query += `$5 = $5 and `
		}
		if len(postgresTags) != 0 {
			query += `tag && $6::varchar[] `
		} else {
			query += `$6 = $6 `
		}
		query += `group by e.id,
         e.title,
         e.description,
         e.text,
         e.city,
         e.category,
         e.viewed,
         e.img_url,
         e.date,
         e.geo,
         e.address,
         e.tag,
         e.author_id 
         order by viewed DESC`

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(query).
			WithArgs(userIdInt, test.title, test.category, test.city, test.date, postgresTags).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)
		out, actualErr := repositoryTest.GetEvents(test.userId, test.title, test.category, test.city, test.date, test.tags)
		require.Equal(t, test.outputErr, actualErr)
		if test.outputErr != nil {
			require.Equal(t, []*models.Event(nil), out)
		} else {
			require.Equal(t, test.outputRes, out)
		}
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
		error2.ErrAtoi,
		[]*models.Event{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error2.ErrPostgres,
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
		out, actualErr := repositoryTest.GetVisitedEvents(test.userId)
		require.Equal(t, test.outputErr, actualErr)
		if test.outputErr != nil {
			require.Equal(t, []*models.Event(nil), out)
		} else {
			require.Equal(t, test.outputRes, out)
		}
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
		error2.ErrAtoi,
		[]*models.Event{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error2.ErrPostgres,
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
		out, actualErr := repositoryTest.GetCreatedEvents(test.userId)
		require.Equal(t, test.outputErr, actualErr)
		if test.outputErr != nil {
			require.Equal(t, []*models.Event(nil), out)
		} else {
			require.Equal(t, test.outputRes, out)
		}
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
		error2.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		nil,
		error2.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		sql2.ErrConnDone,
		error2.ErrPostgres,
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
		actualErr := repositoryTest.Visit(test.eventId, test.userId)
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
		error2.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		nil,
		error2.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		sql2.ErrConnDone,
		error2.ErrPostgres,
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
		actualErr := repositoryTest.Unvisit(test.eventId, test.userId)
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
		error2.ErrAtoi,
	},
	{
		3,
		"1",
		"b",
		0,
		false,
		nil,
		error2.ErrAtoi,
	},
	{
		4,
		"1",
		"2",
		0,
		false,
		sql2.ErrConnDone,
		error2.ErrPostgres,
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
		out, actualErr := repositoryTest.IsVisited(test.eventId, test.userId)
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.result, out)
	}
}

var getCitiesTests = []struct {
	id           int
	postgresErr  error
	outputErr    error
	outputResult []string
}{
	{
		1,
		nil,
		nil,
		[]string{"test"},
	},
}

func TestGetCities(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getCitiesTests {

		rows := sqlmock.NewRows([]string{"city"}).AddRow("test")

		mock.ExpectQuery(getCitiesQuery).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)
		out, actualErr := repositoryTest.GetCities()
		require.Equal(t, test.outputErr, actualErr)
		require.Equal(t, test.outputResult, out)
	}
}
