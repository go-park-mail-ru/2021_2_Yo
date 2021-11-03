package http

import (
	"backend/models"
	"backend/response"
	image "backend/service/image/manager"
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
		error3.ErrUserNotFound,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
}

func TestGetUser(t *testing.T) {
	for _, test := range getUserTests {
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)

		userId := test.input
		useCaseMock.On("GetUser", userId).Return(test.user, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/user", handlerTest.GetUser).Methods("GET")
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/user", bytes.NewBuffer(nil))
		require.NoError(t, err, logTestMessage+"NewRequest error")
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
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
		error3.ErrUserNotFound,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
}

func TestGetUserById(t *testing.T) {
	for _, test := range getUserByIdTests {
		userId := test.input
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)
		useCaseMock.On("GetUser", userId).Return(test.user, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/user", handlerTest.GetUserById).Methods("GET")
		req, err := http.NewRequest("GET", "/user", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		//TODO: Баг библиотеки, из-за него пришлось делать дополнительную миддлвару, достающую значения из контекста
		/*
			req = mux.SetURLVars(req, map[string]string{
				"id": userId,
			})
		*/
		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "id", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
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
	output     *response.Response
}{
	{1,
		"1",
		&models.ResponseBodyUser{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		&models.ResponseBodyUser{
			Name:    "testName",
			Surname: "testSurname",
			About:   "testAbout",
		},
		error3.ErrUserNotFound,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
}

func TestUpdateUserInfo(t *testing.T) {
	for _, test := range updateUserInfoTests {
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)

		userId := test.input

		userModel := new(models.User)
		if test.user != nil {
			userModel.Name = test.user.Name
			userModel.Surname = test.user.Surname
			userModel.About = test.user.About
		}

		useCaseMock.On("UpdateUserInfo",
			userId,
			userModel.Name,
			userModel.Surname,
			userModel.About).Return(test.useCaseErr)

		bodyUserJSON, err := json.Marshal(test.user)
		require.NoError(t, err, logTestMessage+"err =", err)

		if test.user == nil {
			bodyUserJSON = nil
		}

		r := mux.NewRouter()
		r.HandleFunc("/user/info", handlerTest.UpdateUserInfo).Methods("POST")
		req, err := http.NewRequest("POST", "/user/info", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
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
	output     *response.Response
}{
	{1,
		"1",
		&models.ResponseBodyUser{
			Password: "testPassword",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		&models.ResponseBodyUser{
			Password: "testPassword",
		},
		error3.ErrUserNotFound,
		response.ErrorResponse(error3.ErrUserNotFound.Error())},
	{3,
		"1",
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
}

func TestUpdateUserPassword(t *testing.T) {
	for _, test := range updateUserPasswordTests {
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)

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
		r.HandleFunc("/user/password", handlerTest.UpdateUserPassword).Methods("POST")
		req, err := http.NewRequest("POST", "/user/password", bytes.NewBuffer(bodyUserJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

//TODO: TestUpdateUserPhoto

var updateUserAvatarTests = []struct {
	id     int
	input  string
	output *response.Response
}{
	{1,
		"1",
		response.OkResponse()},
}

func TestUpdateUserAvatar(t *testing.T) {
	for _, test := range updateUserAvatarTests {
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)

		userId := test.input

		r := mux.NewRouter()
		r.HandleFunc("/user/avatar", handlerTest.UpdateUserAvatar).Methods("POST")
		req, err := http.NewRequest("POST", "/user/avatar", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}
