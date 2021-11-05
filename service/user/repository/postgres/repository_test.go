package postgres

import (
	"backend/models"
	error3 "backend/service/user/error"
	sql2 "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:auth:repository:postgres:"

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

		actualUser, actualErr := repositoryTest.GetUserById(test.userId)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		if test.outputErr == nil {
			require.Equal(t, test.outputUser, actualUser, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		} else {
			require.Equal(t, (*models.User)(nil), actualUser, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		}
	}
}

var updateUserInfoTests = []struct {
	id          int
	userId      string
	name        string
	surname     string
	about       string
	postgresErr error
	outputErr   error
}{
	{
		1,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		nil,
		nil,
	},
	{
		2,
		"a",
		"testName",
		"testSurname",
		"testAbout",
		nil,
		error3.ErrAtoi,
	},
	{
		3,
		"1",
		"testName",
		"testSurname",
		"testAbout",
		sql2.ErrConnDone,
		error3.ErrPostgres,
	},
}

func TestUpdateUserInfo(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range updateUserInfoTests {
		userIdInt, err := strconv.Atoi(test.userId)
		if err != nil {
			userIdInt = 0
		}
		mock.ExpectQuery(updateUserInfoQuery).
			WithArgs(test.name, test.surname, test.about, userIdInt).
			WillReturnRows(sqlmock.NewRows([]string{})).
			WillReturnError(test.postgresErr)

		actualErr := repositoryTest.UpdateUserInfo(test.userId, test.name, test.surname, test.about)
		t.Log(actualErr)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
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

		actualErr := repositoryTest.UpdateUserPassword(test.userId, test.password)
		t.Log(actualErr)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
