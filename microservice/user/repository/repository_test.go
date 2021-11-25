package repository

import (
	userGrpc "backend/microservice/user/proto"
	"backend/pkg/models"
	error3 "backend/service/user/error"
	"context"
	sql2 "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func fromProtoToModel(u *userGrpc.User) *models.User {
	return &models.User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Mail:     u.Mail,
		Password: u.Password,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
	}
}

var getUserByIdTests = []struct {
	id          int
	userId      string
	postgresErr error
	outputUser  *models.User
	outputErr   error
}{
	{
		1,
		"1",
		nil,
		&models.User{
			ID:       "1",
			Name:     "testName",
			Surname:  "testSurname",
			Mail:     "testMail",
			Password: "testPassword",
			About:    "testAbout",
		},
		nil,
	},
	{
		2,
		"1",
		sql2.ErrNoRows,
		&models.User{},
		error3.ErrUserNotFound,
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		&models.User{},
		error3.ErrPostgres,
	},
}

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getUserByIdTests {
		mock.ExpectQuery(getUserByIdQuery).WithArgs(test.userId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "mail", "password", "about"}).
				AddRow(test.outputUser.ID,
					test.outputUser.Name,
					test.outputUser.Surname,
					test.outputUser.Mail,
					test.outputUser.Password,
					test.outputUser.About)).WillReturnError(test.postgresErr)

		in := &userGrpc.UserId{ID: test.userId}
		out, actualErr := repositoryTest.GetUserById(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualUser := fromProtoToModel(out)
		require.Equal(t, test.outputUser, actualUser)
	}
}

var updateUserInfoTests = []struct {
	id          int
	userId      string
	name        string
	surname     string
	about       string
	imgUrl      string
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		"",
		nil,
		nil,
	},
	{
		2,
		"a",
		"testName",
		"testSurname",
		"testAbout",
		"",
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		"",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
	{
		4,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		"img",
		nil,
		nil,
	},
	{
		4,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		"img",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err, logMessage, err)
		defer db.Close()
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repositoryTest := NewRepository(sqlxDB)
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}
		if test.imgUrl != "" {
			mock.ExpectQuery(updateUserInfoQuery).
				WithArgs(test.name, test.surname, test.about, test.imgUrl, userIdInt).
				WillReturnRows(sqlmock.NewRows([]string{})).
				WillReturnError(test.postgresErr)
		} else {
			mock.ExpectQuery(updateUserInfoQueryWithoutImgUrl).
				WithArgs(test.name, test.surname, test.about, userIdInt).
				WillReturnRows(sqlmock.NewRows([]string{})).
				WillReturnError(test.postgresErr)
		}
		in := &userGrpc.User{
			ID:      test.userId,
			Name:    test.name,
			Surname: test.surname,
			About:   test.about,
			ImgUrl:  test.imgUrl,
		}
		_, actualErr := repositoryTest.UpdateUserInfo(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var updateUserPasswordTests = []struct {
	id          int
	userId      string
	password    string
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"testPassword",
		nil,
		nil,
	},
	{
		2,
		"a",
		"testPassword",
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"testPassword",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range updateUserPasswordTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}
		mock.ExpectQuery(updateUserPasswordQuery).
			WithArgs(test.password, userIdInt).
			WillReturnRows(sqlmock.NewRows([]string{})).
			WillReturnError(test.postgresErr)

		in := &userGrpc.UpdateUserPasswordRequest{
			ID:       test.userId,
			Password: test.password,
		}
		_, actualErr := repositoryTest.UpdateUserPassword(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var getSubscribersTests = []struct {
	id          int
	userId      string
	postgresErr error
	outputErr   error
	outputRes   []*models.User
}{
	{
		1,
		"1",
		nil,
		nil,
		[]*models.User{
			&models.User{ID: "1"},
		},
	},
	{
		2,
		"a",
		nil,
		error3.ErrAtoi,
		[]*models.User{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error3.ErrPostgres,
		[]*models.User{},
	},
}

func TestGetSubscribers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getSubscribersTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(getSubscribersQuery).
			WithArgs(userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.UserId{
			ID: test.userId,
		}
		out, actualErr := repositoryTest.GetSubscribers(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.User, len(out.Users))
		for i, u := range out.Users {
			actualRes[i] = fromProtoToModel(u)
		}
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getSubscribesTests = []struct {
	id          int
	userId      string
	postgresErr error
	outputErr   error
	outputRes   []*models.User
}{
	{
		1,
		"1",
		nil,
		nil,
		[]*models.User{
			&models.User{ID: "1"},
		},
	},
	{
		2,
		"a",
		nil,
		error3.ErrAtoi,
		[]*models.User{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error3.ErrPostgres,
		[]*models.User{},
	},
}

func TestGetSubscribes(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getSubscribesTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(getSubscribesQuery).
			WithArgs(userIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.UserId{
			ID: test.userId,
		}
		out, actualErr := repositoryTest.GetSubscribes(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.User, len(out.Users))
		for i, u := range out.Users {
			actualRes[i] = fromProtoToModel(u)
		}
		require.Equal(t, test.outputRes, actualRes)
	}
}

var getVisitorsTests = []struct {
	id          int
	eventId     string
	postgresErr error
	outputErr   error
	outputRes   []*models.User
}{
	{
		1,
		"1",
		nil,
		nil,
		[]*models.User{
			&models.User{ID: "1"},
		},
	},
	{
		2,
		"a",
		nil,
		error3.ErrAtoi,
		[]*models.User{},
	},
	{
		3,
		"1",
		sql2.ErrConnDone,
		error3.ErrPostgres,
		[]*models.User{},
	},
}

func TestGetVisitors(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getVisitorsTests {
		eventIdInt, err := strconv.Atoi(test.eventId)
		if err != nil {
			eventIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(getVisitorsQuery).
			WithArgs(eventIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.EventId{
			ID: test.eventId,
		}
		out, actualErr := repositoryTest.GetVisitors(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := make([]*models.User, len(out.Users))
		for i, u := range out.Users {
			actualRes[i] = fromProtoToModel(u)
		}
		require.Equal(t, test.outputRes, actualRes)
	}
}

var subscribeTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	postgresErr  error
	outputErr    error
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

func TestSubscribe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range subscribeTests {
		subscribedIdInt, err := strconv.Atoi(test.subscribedId)
		if err != nil {
			subscribedIdInt = 0
		}
		subscriberIdInt, err := strconv.Atoi(test.subscriberId)
		if err != nil {
			subscriberIdInt = 0
		}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(subscribeQuery).
			WithArgs(subscribedIdInt, subscriberIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		_, actualErr := repositoryTest.Subscribe(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var unsubscribeTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	postgresErr  error
	outputErr    error
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

func TestUnsubscribe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range unsubscribeTests {
		subscribedIdInt, err := strconv.Atoi(test.subscribedId)
		if err != nil {
			subscribedIdInt = 0
		}
		subscriberIdInt, err := strconv.Atoi(test.subscriberId)
		if err != nil {
			subscriberIdInt = 0
		}

		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(unsubscribeQuery).
			WithArgs(subscribedIdInt, subscriberIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		_, actualErr := repositoryTest.Unsubscribe(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
	}
}

var isSubscribedTests = []struct {
	id           int
	subscribedId string
	subscriberId string
	count        int
	result       bool
	postgresErr  error
	outputErr    error
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

func TestIsSubscribed(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range isSubscribedTests {
		subscribedIdInt, err := strconv.Atoi(test.subscribedId)
		if err != nil {
			subscribedIdInt = 0
		}
		subscriberIdInt, err := strconv.Atoi(test.subscriberId)
		if err != nil {
			subscriberIdInt = 0
		}

		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(test.count)
		mock.ExpectQuery(isSubscribedQuery).
			WithArgs(subscribedIdInt, subscriberIdInt).
			WillReturnRows(rows).
			WillReturnError(test.postgresErr)

		in := &userGrpc.SubscribeRequest{
			SubscribedId: test.subscribedId,
			SubscriberId: test.subscriberId,
		}
		out, actualErr := repositoryTest.IsSubscribed(context.Background(), in)
		require.Equal(t, test.outputErr, actualErr)
		actualRes := out.Result
		require.Equal(t, test.result, actualRes)
	}
}
