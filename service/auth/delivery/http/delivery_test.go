package http

import (
	"backend/pkg/models"
	"backend/pkg/response"
	error3 "backend/service/auth/error"
	"backend/service/auth/usecase"
	error4 "backend/service/csrf/error"
	csrf "backend/service/csrf/manager"
	error2 "backend/service/session/error"
	session "backend/service/session/manager"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const logTestMessage = "auth:delivery:test"

var signUpTests = []struct {
	id                int
	input             *models.ResponseBodyUser
	useCaseErr        error
	sessionManagerErr error
	csrfManagerErr    error
	output            *response.Response
}{
	{1,
		&models.ResponseBodyUser{
			Mail: "testMail@mail.ru",
		},
		nil,
		nil,
		nil,
		response.OkResponse()},
	{2,
		&models.ResponseBodyUser{
			Mail: "testMail",
		},
		nil,
		nil,
		nil,
		response.ErrorResponse(response.ErrValidation.Error())},
	{3,
		&models.ResponseBodyUser{
			Mail: "testMail@mail.ru",
		},
		error3.ErrUserNotFound,
		nil,
		nil,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
	{4,
		&models.ResponseBodyUser{
			Mail: "testMail@mail.ru",
		},
		nil,
		error2.ErrEmptySessionId,
		nil,
		response.ErrorResponse(error2.ErrEmptySessionId.Error())},
	{5,
		&models.ResponseBodyUser{
			Mail: "testMail@mail.ru",
		},
		nil,
		nil,
		error4.ErrEmptyToken,
		response.ErrorResponse(error4.ErrEmptyToken.Error())},
}

func TestSignUp(t *testing.T) {
	for _, test := range signUpTests {
		useCaseMock := new(usecase.UseCaseMock)
		sessionManagerMock := new(session.ManagerMock)
		csrfManagerMock := new(csrf.ManagerMock)
		deliveryTest := NewDelivery(useCaseMock, sessionManagerMock, csrfManagerMock)

		bodyUserJSON, err := json.Marshal(test.input)
		require.NoError(t, err, logTestMessage+"err =", err)

		userModel := new(models.User)
		userModel.Name = test.input.Name
		userModel.Surname = test.input.Surname
		userModel.About = test.input.About
		userModel.Mail = test.input.Mail
		userModel.Password = test.input.Password

		useCaseMock.On("SignUp", userModel).Return("", test.useCaseErr)
		sessionManagerMock.On("Create", "").Return("", test.sessionManagerErr)
		csrfManagerMock.On("Create", "").Return("", test.csrfManagerErr)

		r := mux.NewRouter()
		r.HandleFunc("/signup", deliveryTest.SignUp).Methods("POST")
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var signInTests = []struct {
	id                int
	input             *models.ResponseBodyUser
	useCaseErr        error
	sessionManagerErr error
	csrfManagerErr    error
	output            *response.Response
}{
	{1,
		&models.ResponseBodyUser{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		nil,
		nil,
		response.OkResponse()},
	{2,
		&models.ResponseBodyUser{
			Mail:     "testMail",
			Password: "testPassword",
		},
		nil,
		nil,
		nil,
		response.ErrorResponse(response.ErrValidation.Error())},
	{3,
		&models.ResponseBodyUser{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		error3.ErrUserNotFound,
		nil,
		nil,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
	{4,
		&models.ResponseBodyUser{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		error2.ErrEmptySessionId,
		nil,
		response.ErrorResponse(error2.ErrEmptySessionId.Error())},
	{5,
		&models.ResponseBodyUser{
			Mail:     "testMail@mail.ru",
			Password: "testPassword",
		},
		nil,
		nil,
		error4.ErrEmptyToken,
		response.ErrorResponse(error4.ErrEmptyToken.Error())},
}

func TestSignIn(t *testing.T) {
	for _, test := range signInTests {
		useCaseMock := new(usecase.UseCaseMock)
		sessionManagerMock := new(session.ManagerMock)
		csrfManagerMock := new(csrf.ManagerMock)
		deliveryTest := NewDelivery(useCaseMock, sessionManagerMock, csrfManagerMock)

		userMail := test.input.Mail
		userPassword := test.input.Password

		useCaseMock.On("SignIn", userMail, userPassword).Return("", test.useCaseErr)
		sessionManagerMock.On("Create", "").Return("", test.sessionManagerErr)
		csrfManagerMock.On("Create", "").Return("", test.csrfManagerErr)

		bodyUserJSON, err := json.Marshal(test.input)
		require.NoError(t, err, logTestMessage+"err =", err)

		r := mux.NewRouter()
		r.HandleFunc("/login", deliveryTest.SignIn).Methods("POST")
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var logoutTests = []struct {
	id                int
	input             *http.Cookie
	csrfToken         string
	sessionManagerErr error
	csrfManagerErr    error
	output            *response.Response
}{
	{1,
		&http.Cookie{
			Name:  "session_id",
			Value: "123",
		},
		"token",
		nil,
		nil,
		response.OkResponse()},
	{2,
		&http.Cookie{
			Name:  "",
			Value: "",
		},
		"token",
		nil,
		nil,
		response.ErrorResponse(error3.ErrCookie.Error())},
	{3,
		&http.Cookie{
			Name:  "session_id",
			Value: "",
		},
		"token",
		error2.ErrDeleteSession,
		nil,
		response.ErrorResponse(error2.ErrDeleteSession.Error())},
	{4,
		&http.Cookie{
			Name:  "session_id",
			Value: "",
		},
		"token",
		nil,
		error4.ErrRedis,
		response.ErrorResponse(error4.ErrRedis.Error())},
}

func TestLogout(t *testing.T) {
	for _, test := range logoutTests {
		useCaseMock := new(usecase.UseCaseMock)
		sessionManagerMock := new(session.ManagerMock)
		csrfManagerMock := new(csrf.ManagerMock)
		deliveryTest := NewDelivery(useCaseMock, sessionManagerMock, csrfManagerMock)

		cookie := test.input
		csrfToken := test.csrfToken

		sessionManagerMock.On("Delete", cookie.Value).Return(test.sessionManagerErr)
		csrfManagerMock.On("Delete", csrfToken).Return(test.csrfManagerErr)

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

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
