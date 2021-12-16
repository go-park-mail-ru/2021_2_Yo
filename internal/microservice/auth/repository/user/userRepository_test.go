package user

import (
	"backend/internal/models"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		errors.New("internal DB server error"),
	},
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err, logMessage, err)
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
		assert.Equal(t, test.outputErr, actualErr)
		assert.Equal(t, test.outputStr, actualStr)
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
		errors.New("internal DB server error"),
		&models.User{
			ID:       "1",
			Name:     "testName",
			Surname:  "testSurname",
			Mail:     "testMail",
			Password: "testPassword",
			About:    "testAbout",
		},
		errors.New("internal DB server error"),
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
		errors.New("internal DB server error"),
	},
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err, logMessage, err)
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
		assert.Equal(t, test.outputErr, actualErr)
		var expectedUser *models.User
		if test.outputErr != nil {
			expectedUser = nil
		} else {
			expectedUser = test.outputUser
		}
		assert.Equal(t, expectedUser, actualUser)
	}
}
