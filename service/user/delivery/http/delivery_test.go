package http

import (
	"backend/models"
	"backend/response"
	image "backend/service/image/manager"
	error3 "backend/service/user/error"
	"backend/service/user/usecase"
	"bytes"
	"context"
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
		userId := test.input
		useCaseMock := new(usecase.UseCaseMock)
		imageManagerMock := new(image.ManagerMock)
		handlerTest := NewDelivery(useCaseMock, imageManagerMock)
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
