package http

import (
	"backend/auth/usecase"
	"backend/models"
	"backend/response"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	
	r := mux.NewRouter()
	useCaseMock := new(usecase.UseCaseAuthMock)
	handlerTest := NewHandlerAuth(useCaseMock)
	r.HandleFunc("/signup", handlerTest.SignUp).Methods("POST")

	bodyUserTest := &response.ResponseBodyUser{
		Name:     "nameTest",
		Surname:  "surnameTest",
		Mail:     "mailTest",
		Password: "passwordTest",
	}

	bodyUserJSON, err := json.Marshal(bodyUserTest)
	require.NoError(t, err, "TestSignUp : jsonMarshal error = ", err)

	useCaseMock.On(
		"SignUp",
		bodyUserTest.Name,
		bodyUserTest.Surname,
		bodyUserTest.Mail,
		bodyUserTest.Password).Return(nil)

	useCaseMock.On(
		"SignIn",
		bodyUserTest.Mail,
		bodyUserTest.Password).Return("", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(bodyUserJSON))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, "TestSignUp : expected 200, got", w.Code)
}

func TestSignIn(t *testing.T) {
	r := mux.NewRouter()
	useCaseMock := new(usecase.UseCaseAuthMock)
	handlerTest := NewHandlerAuth(useCaseMock)
	r.HandleFunc("/login", handlerTest.SignIn).Methods("POST")

	bodyUserTest := &response.ResponseBodyUser{
		Mail:     "mailTest",
		Password: "passwordTest",
	}

	bodyUserJSON, err := json.Marshal(bodyUserTest)
	require.NoError(t, err, "TestSignIn : jsonMarshal error = ", err)

	useCaseMock.On(
		"SignIn",
		bodyUserTest.Mail,
		bodyUserTest.Password).Return("jwt_token_test", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(bodyUserJSON))
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, "TestSignIn : expected 200, got", w.Code)
	require.Equal(t, "jwt_token_test", w.Result().Cookies()[0].Value, "TestSignIn : expected jwt_token_test, got", w.Result().Cookies()[0].Value)
}

func TestUser(t *testing.T) {
	r := mux.NewRouter()
	useCaseMock := new(usecase.UseCaseAuthMock)
	handlerTest := NewHandlerAuth(useCaseMock)
	r.HandleFunc("/user", handlerTest.User).Methods("GET")

	w := httptest.NewRecorder()
	jwtToken := "test_token"
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: jwtToken,
	}

	useCaseMock.On(
		"ParseToken",
		cookie.Value).Return(&models.User{
		ID:       "1",
		Name:     "nameTest",
		Surname:  "surnameTest",
		Mail:     "mailTest",
		Password: "passwordTest",
	}, nil)

	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer([]byte("")))
	req.AddCookie(cookie)
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code, "TestSignIn : expected 200, got", w.Code)
	require.Equal(t,
		"{\"status\":200,\"body\":{\"name\":\"nameTest\"}}",
		w.Body.String(),
		"TestSignIn : expected jwt_token_test, got",
		w.Body.String())
}
