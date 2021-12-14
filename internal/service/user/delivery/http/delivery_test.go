package http

import (
	models "backend/internal/models"
	response "backend/internal/response"
	error2 "backend/internal/service/user/error"
	"backend/internal/service/user/usecase"
	"backend/pkg/response"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

const logTestMessage = "auth:delivery:test"

var getUserTests = []struct {
	id         int
	input      string
	user       *models.User
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		&models.User{
			ID: "1",
		},
		nil,
		response.UserResponse(&models.User{ID: "1"})},
	{2,
		"1",
		nil,
		error2.ErrUserNotFound,
		response.ErrorResponse(error2.ErrUserNotFound.Error())},
}

func TestGetUser(t *testing.T) {
	for _, test := range getUserTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input
		useCaseMock.On("GetUserById", userId).Return(test.user, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/user", deliveryTest.GetUser).Methods("GET")
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/user", bytes.NewBuffer(nil))
		require.NoError(t, err, logTestMessage+"NewRequest error")
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))
	}
}

var getUserByIdTests = []struct {
	id         int
	input      string
	user       *models.User
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		&models.User{
			ID: "1",
		},
		nil,
		response.UserResponse(&models.User{
			ID: "1",
		})},
	{2,
		"1",
		nil,
		error2.ErrUserNotFound,
		response.ErrorResponse(error2.ErrUserNotFound.Error())},
}

func TestGetUserById(t *testing.T) {
	for _, test := range getUserByIdTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input

		useCaseMock.On("GetUserById", userId).Return(test.user, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/user/{id:[0-9]+}", deliveryTest.GetUserById).Methods("GET")
		req, err := http.NewRequest("GET", "/user/"+userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var updateUserInfoTests = []struct {
	id         int
	input      string
	user       *models.UserResponseBody
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		&models.UserResponseBody{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		&models.UserResponseBody{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		error2.ErrUserNotFound,
		response.ErrorResponse(error2.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input

		userModel := new(models.User)
		if test.user != nil {
			userModel.ID = userId
			userModel.Name = test.user.Name
			userModel.Surname = test.user.Surname
			userModel.About = test.user.About
		}

		useCaseMock.On("UpdateUserInfo", userModel).Return(test.useCaseErr)

		bodyUserJSON, err := json.Marshal(test.user)
		require.NoError(t, err, logTestMessage+"err =", err)

		if test.user == nil {
			bodyUserJSON = nil
		}

		r := mux.NewRouter()
		r.HandleFunc("/user/info", deliveryTest.UpdateUserInfo).Methods("POST")
		req, err := http.NewRequest("POST", "/user/info", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))
	}
}

var updateUserPasswordTests = []struct {
	id         int
	input      string
	user       *models.UserResponseBody
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		&models.UserResponseBody{
			Password: "testPassword",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		&models.UserResponseBody{
			Password: "testPassword",
		},
		error2.ErrUserNotFound,
		response.ErrorResponse(error2.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
}

func TestUpdateUserPassword(t *testing.T) {
	for _, test := range updateUserPasswordTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input

		userModel := new(models.User)
		if test.user != nil {
			userModel.Password = test.user.Password
		}

		useCaseMock.On("UpdateUserPassword",
			userId,
			userModel.Password).Return(test.useCaseErr)

		bodyUserJSON, err := json.Marshal(test.user)
		require.NoError(t, err, logTestMessage+"err =", err)

		if test.user == nil {
			bodyUserJSON = nil
		}

		r := mux.NewRouter()
		r.HandleFunc("/user/password", deliveryTest.UpdateUserPassword).Methods("POST")
		req, err := http.NewRequest("POST", "/user/password", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))
	}
}

var getSubscribersTests = []struct {
	id         int
	userId     string
	useCaseErr error
}{
	{
		1,
		"1",
		nil,
	},
	{
		2,
		"1",
		errors.New("test_err"),
	},
}

func TestGetSubscribers(t *testing.T) {
	for _, test := range getSubscribersTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetSubscribers", test.userId).Return([]*models.User{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/{id:[0-9]+}", deliveryTest.GetSubscribers).Methods("GET")
		req, err := http.NewRequest("GET", "/"+test.userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getSubscribesTests = []struct {
	id         int
	userId     string
	useCaseErr error
}{
	{
		1,
		"1",
		nil,
	},
	{
		2,
		"1",
		errors.New("test_err"),
	},
}

func TestGetSubscribes(t *testing.T) {
	for _, test := range getSubscribesTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetSubscribes", test.userId).Return([]*models.User{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/{id:[0-9]+}", deliveryTest.GetSubscribes).Methods("GET")
		req, err := http.NewRequest("GET", "/"+test.userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getVisitorsTests = []struct {
	id         int
	eventId    string
	useCaseErr error
}{
	{
		1,
		"1",
		nil,
	},
	{
		2,
		"1",
		errors.New("test_err"),
	},
}

func TestGetVisitors(t *testing.T) {
	for _, test := range getVisitorsTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetVisitors", test.eventId).Return([]*models.User{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.GetVisitors).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		eventIdCtx := context.WithValue(context.Background(), "eventId", test.eventId)
		r.ServeHTTP(w, req.WithContext(eventIdCtx))
	}
}

var subscribeTests = []struct {
	id         int
	vars       interface{}
	userId     interface{}
	useCaseErr error
}{
	{
		1,
		map[string]string{
			"id": "123",
		},
		"1",
		nil,
	},
	{
		2,
		map[string]string{
			"id": "123",
		},
		"1",
		errors.New("test_err"),
	},
}

func TestSubscribe(t *testing.T) {
	for _, test := range subscribeTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		var eId string
		var uId string
		vars, ok := test.vars.(map[string]string)
		if ok {
			eId = vars["id"]
		}
		userId, ok := test.userId.(string)
		if ok {
			uId = userId
		}

		useCaseMock.On("Subscribe", eId, uId).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.Subscribe).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var unsubscribeTests = []struct {
	id         int
	vars       interface{}
	userId     interface{}
	useCaseErr error
}{
	{
		1,
		map[string]string{
			"id": "123",
		},
		"1",
		nil,
	},
	{
		2,
		map[string]string{
			"id": "123",
		},
		"1",
		errors.New("test_err"),
	},
}

func TestUnsubscribe(t *testing.T) {
	for _, test := range unsubscribeTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		var eId string
		var uId string
		vars, ok := test.vars.(map[string]string)
		if ok {
			eId = vars["id"]
		}
		userId, ok := test.userId.(string)
		if ok {
			uId = userId
		}

		useCaseMock.On("Unsubscribe", eId, uId).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.Unsubscribe).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var isSubscribedTests = []struct {
	id         int
	vars       interface{}
	userId     interface{}
	useCaseErr error
}{
	{
		1,
		map[string]string{
			"id": "123",
		},
		"1",
		nil,
	},
	{
		2,
		map[string]string{
			"id": "123",
		},
		"1",
		errors.New("test_err"),
	},
}

func TestIsSubscribed(t *testing.T) {
	for _, test := range isSubscribedTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		var eId string
		var uId string
		vars, ok := test.vars.(map[string]string)
		if ok {
			eId = vars["id"]
		}
		userId, ok := test.userId.(string)
		if ok {
			uId = userId
		}

		useCaseMock.On("IsSubscribed", eId, uId).Return(false, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.IsSubscribed).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
