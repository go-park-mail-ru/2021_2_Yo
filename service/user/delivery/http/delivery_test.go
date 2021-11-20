package http

import (
	"backend/models"
	models2 "backend/pkg/models"
	response2 "backend/pkg/response"
	error3 "backend/service/user/error"
	"backend/service/user/usecase"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const logTestMessage = "auth:delivery:test"

var getUserTests = []struct {
	id         int
	input      string
	user       *models2.User
	useCaseErr error
	output     *response2.Response
}{
	{1,
		"1",
		&models2.User{
			ID: "1",
		},
		nil,
		response2.UserResponse(&models2.User{ID: "1"})},
	{2,
		"1",
		nil,
		error3.ErrUserNotFound,
		response2.ErrorResponse(error3.ErrUserNotFound.Error())},
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

		wTest := httptest.NewRecorder()
		response2.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getUserByIdTests = []struct {
	id         int
	input      string
	user       *models2.User
	useCaseErr error
	output     *response2.Response
}{
	{1,
		"1",
		&models2.User{
			ID: "1",
		},
		nil,
		response2.UserResponse(&models2.User{
			ID: "1",
		})},
	{2,
		"1",
		nil,
		error3.ErrUserNotFound,
		response2.ErrorResponse(error3.ErrUserNotFound.Error())},
}

func TestGetUserById(t *testing.T) {
	for _, test := range getUserByIdTests {
		userId := test.input
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)
		useCaseMock.On("GetUserById", userId).Return(test.user, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/user/{id:[0-9]+}", deliveryTest.GetUserById).Methods("GET")
		req, err := http.NewRequest("GET", "/user/"+userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		wTest := httptest.NewRecorder()
		response2.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateUserInfoTests = []struct {
	id         int
	input      string
	user       *models.ResponseBodyUser
	useCaseErr error
	output     *response2.Response
}{
	{1,
		"1",
		&models.ResponseBodyUser{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		nil,
		response2.OkResponse()},
	{2,
		"1",
		&models.ResponseBodyUser{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		error3.ErrUserNotFound,
		response2.ErrorResponse(error3.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response2.ErrorResponse(response2.ErrJSONDecoding.Error())},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input

		userModel := new(models2.User)
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

		wTest := httptest.NewRecorder()
		response2.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateUserPasswordTests = []struct {
	id         int
	input      string
	user       *models.ResponseBodyUser
	useCaseErr error
	output     *response2.Response
}{
	{1,
		"1",
		&models.ResponseBodyUser{
			Password: "testPassword",
		},
		nil,
		response2.OkResponse()},
	{2,
		"1",
		&models.ResponseBodyUser{
			Password: "testPassword",
		},
		error3.ErrUserNotFound,
		response2.ErrorResponse(error3.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response2.ErrorResponse(response2.ErrJSONDecoding.Error())},
}

func TestUpdateUserPassword(t *testing.T) {
	for _, test := range updateUserPasswordTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		userId := test.input

		userModel := new(models2.User)
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

		wTest := httptest.NewRecorder()
		response2.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
