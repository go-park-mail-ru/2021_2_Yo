package postgres

import (
	"backend/models"
	error2 "backend/service/auth/error"
	sql2 "database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const logTestMessage = "service:auth:repository:postgres:"

var createUserTests = []struct {
	id        int
	user      *models.User
	userId    string
	outputStr string
	outputErr error
}{
	{
		1,
		&models.User{
			Name:     "testName",
			Surname:  "testSurname",
			Mail:     "testMail",
			Password: "testPassword",
			About:    "testAbout",
		},
		"1",
		"1",
		nil,
	},
	{
		2,
		&models.User{},
		"",
		"",
		error2.ErrPostgres,
	},
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range createUserTests {
		testPostgresUser := toPostgresUser(test.user)

		mock.ExpectQuery(createUserQuery).WithArgs(
			testPostgresUser.Name,
			testPostgresUser.Surname,
			testPostgresUser.Mail,
			testPostgresUser.Password,
			testPostgresUser.About).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.userId)).WillReturnError(test.outputErr)

		actualStr, actualErr := repositoryTest.CreateUser(test.user)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		require.Equal(t, test.outputStr, actualStr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getUserTests = []struct {
	id          int
	mail        string
	password    string
	postgresErr error
	outputUser  *models.User
	outputErr   error
}{
	{
		1,
		"testMail",
		"testPassword",
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
		"testMail",
		"testPassword",
		sql2.ErrNoRows,
		&models.User{
			ID:       "1",
			Name:     "testName",
			Surname:  "testSurname",
			Mail:     "testMail",
			Password: "testPassword",
			About:    "testAbout",
		},
		error2.ErrUserNotFound,
	},
	{
		3,
		"testMail",
		"testPassword",
		errors.New("test error"),
		&models.User{
			ID:       "1",
			Name:     "testName",
			Surname:  "testSurname",
			Mail:     "testMail",
			Password: "testPassword",
			About:    "testAbout",
		},
		error2.ErrPostgres,
	},
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err, logMessage, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repositoryTest := NewRepository(sqlxDB)

	for _, test := range getUserTests {
		mock.ExpectQuery(getUserQuery).WithArgs(
			test.mail,
			test.password).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "mail", "password", "about"}).
				AddRow(test.outputUser.ID,
					test.outputUser.Name,
					test.outputUser.Surname,
					test.outputUser.Mail,
					test.outputUser.Password,
					test.outputUser.About)).
			WillReturnError(test.postgresErr)

		actualUser, actualErr := repositoryTest.GetUser(test.mail, test.password)
		require.Equal(t, test.outputErr, actualErr, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
		var expectedUser *models.User
		if test.outputErr != nil {
			expectedUser = nil
		} else {
			expectedUser = test.outputUser
		}
		require.Equal(t, expectedUser, actualUser, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
