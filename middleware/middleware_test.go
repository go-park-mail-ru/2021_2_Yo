package middleware

import (
	"backend/service/auth/usecase"
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {

}

func testHandlerFuncPanic(w http.ResponseWriter, r *http.Request) {
	panic("test panic")
}

var utilsTests = []struct {
	id          int
	err         error
	handlerFunc func(w http.ResponseWriter, r *http.Request)
}{
	{
		1,
		nil,
		testHandlerFunc,
	},
	{
		2,
		errors.New("test error"),
		testHandlerFunc,
	},
	{
		3,
		nil,
		testHandlerFuncPanic,
	},
}

func TestUtils(t *testing.T) {
	for _, test := range utilsTests {

		useCaseMock := new(usecase.UseCaseMock)
		middlewares := NewMiddlewares(useCaseMock)

		//useCaseMock.On("CreateEvent", eventModel, test.userId).Return(test.eventId, test.useCaseErr)

		r := mux.NewRouter()
		r.Use(middlewares.Recovery)
		r.Use(middlewares.CORS)
		r.Use(middlewares.Logging)
		r.Use(middlewares.GetVars)

		r1 := mux.NewRouter()
		r1.Use(middlewares.CSRF)

		r2 := mux.NewRouter()
		r2.Use(middlewares.Auth)

		r.HandleFunc("/test", test.handlerFunc).Methods("GET")
		r1.HandleFunc("/test", test.handlerFunc).Methods("GET")
		r2.HandleFunc("/test", test.handlerFunc).Methods("GET")

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test", bytes.NewBuffer(nil))
		require.NoError(t, err)

		req.Header.Set("X-CSRF-Token", "test")
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: "test",
		}
		req.AddCookie(cookie)

		useCaseMock.On("CheckSession", cookie.Value).Return("", test.err)
		useCaseMock.On("CheckToken", req.Header.Get("X-CSRF-Token")).Return("", test.err)

		r.ServeHTTP(w, req)
		if test.id != 3 {
			r1.ServeHTTP(w, req)
			r2.ServeHTTP(w, req)
		}
	}
}
