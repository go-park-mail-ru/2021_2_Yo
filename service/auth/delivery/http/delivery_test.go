package http

import (
	log "backend/pkg/logger"
	"backend/pkg/models"
	"backend/service/auth/usecase"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const logTestMessage = "auth:delivery:test"

var signUpTests = []struct {
	id          int
	input       *models.UserResponseBody
	useCaseErr1 error
	useCaseErr2 error
	useCaseErr3 error
}{
	{
		1,
		&models.UserResponseBody{
			Mail: "testMail@mail.ru",
		},
		nil,
		nil,
		nil,
	},
	{
		2,
		&models.UserResponseBody{
			Mail: "testMail",
		},
		nil,
		nil,
		nil,
	},
	{
		3,
		&models.UserResponseBody{
			Mail: "testMail@mail.ru",
		},
		errors.New("test_err"),
		nil,
		nil,
	},
	{
		4,
		&models.UserResponseBody{
			Mail: "testMail@mail.ru",
		},
		nil,
		errors.New("test_err"),
		nil,
	},
	{
		5,
		&models.UserResponseBody{
			Mail: "testMail@mail.ru",
		},
		nil,
		nil,
		errors.New("test_err"),
	},
}

func TestSignUp(t *testing.T) {
	for _, test := range signUpTests {

		log.Init(logrus.DebugLevel)

		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		bodyUserJSON, err := json.Marshal(test.input)
		require.NoError(t, err, logTestMessage+"err =", err)

		userModel := new(models.User)
		userModel.Mail = test.input.Mail

		useCaseMock.On("SignUp", userModel).Return("", test.useCaseErr1)
		useCaseMock.On("CreateSession", "").Return("", test.useCaseErr2)
		useCaseMock.On("CreateToken", "").Return("", test.useCaseErr3)

		r := mux.NewRouter()
		r.HandleFunc("/signup", deliveryTest.SignUp).Methods("POST")
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")
		_ = req

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var signInTests = []struct {
	id          int
	input       *models.UserResponseBody
	useCaseErr1 error
	useCaseErr2 error
	useCaseErr3 error
}{
	{
		1,
		&models.UserResponseBody{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		nil,
		nil,
	},
	{
		2,
		&models.UserResponseBody{
			Mail:     "testMail",
			Password: "testPassword",
		},
		nil,
		nil,
		nil,
	},
	{
		3,
		&models.UserResponseBody{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		errors.New("test_err"),
		nil,
		nil,
	},
	{
		4,
		&models.UserResponseBody{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		errors.New("test_err"),
		nil,
	},
	{
		5,
		&models.UserResponseBody{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		nil,
		errors.New("test_err"),
	},
}

func TestSignIn(t *testing.T) {
	for _, test := range signInTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userModel := new(models.User)
		userModel.Mail = test.input.Mail
		userModel.Password = test.input.Password

		useCaseMock.On("SignIn", userModel).Return("", test.useCaseErr1)
		useCaseMock.On("CreateSession", "").Return("", test.useCaseErr2)
		useCaseMock.On("CreateToken", "").Return("", test.useCaseErr3)

		bodyUserJSON, err := json.Marshal(test.input)
		require.NoError(t, err, logTestMessage+"err =", err)

		r := mux.NewRouter()
		r.HandleFunc("/login", deliveryTest.SignIn).Methods("POST")
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var logoutTests = []struct {
	id         int
	input      *http.Cookie
	csrfToken  string
	useCaseErr error
}{
	{
		1,
		&http.Cookie{
			Name:  "session_id",
			Value: "123",
		},
		"token",
		nil,
	},
	{
		2,
		&http.Cookie{
			Name:  "",
			Value: "",
		},
		"token",
		nil,
	},
	{
		3,
		&http.Cookie{
			Name:  "session_id",
			Value: "",
		},
		"token",
		errors.New("test_err"),
	},
	{
		4,
		&http.Cookie{
			Name:  "session_id",
			Value: "",
		},
		"token",
		nil,
	},
}

func TestLogout(t *testing.T) {
	for _, test := range logoutTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		cookie := test.input
		csrfToken := test.csrfToken

		useCaseMock.On("DeleteSession", cookie.Value).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/logout", deliveryTest.Logout).Methods("GET")

		req, err := http.NewRequest("GET", "/logout", bytes.NewBuffer([]byte(csrfToken)))
		require.NoError(t, err, logTestMessage+"NewRequest error")
		if cookie.Name != "" {
			req.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		w.Header().Set("X-CSRF-Token", csrfToken)
		r.ServeHTTP(w, req)
	}
}
