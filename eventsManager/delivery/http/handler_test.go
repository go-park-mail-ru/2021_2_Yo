package http

import (
	"backend/eventsManager/usecase"
	//"backend/models"
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	r := mux.NewRouter()
	useCaseMock := new(usecase.UseCaseEventsManagerMock)
	handlerTest := NewHandlerEventsManager(useCaseMock)
	r.HandleFunc("/events", handlerTest.List).Methods("GET")

	//var expected []*models.Event
	//expected = nil
	useCaseMock.On("List")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", bytes.NewBuffer(nil))
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code, "TestSignUp : expected 200, got", w.Code)
}
