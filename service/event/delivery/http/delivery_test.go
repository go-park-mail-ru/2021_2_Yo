package http

import (
	"backend/pkg/models"
	"backend/pkg/response"
	error2 "backend/service/event/error"
	"backend/service/event/usecase"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const logTestMessage = "service:event:delivery"

var createEventTests = []struct {
	id         int
	userId     string
	eventId    string
	event      *models.Event
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		"100",
		&models.Event{
			ID: "1",
		},
		nil,
		response.EventIdResponse("100")},
	{2,
		"1",
		"100",
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
	{3,
		"1",
		"100",
		&models.Event{
			ID: "1",
		},
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestCreateEvent(t *testing.T) {
	for _, test := range createEventTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		eventModel := new(models.Event)
		if test.event != nil {
			eventModel.ID = test.event.ID
		}

		useCaseMock.On("CreateEvent", eventModel, test.userId).Return(test.eventId, test.useCaseErr)

		bodyEventJSON, err := json.Marshal(test.event)
		require.NoError(t, err, logTestMessage+"err =", err)
		if test.event == nil {
			bodyEventJSON = nil
		}

		r := mux.NewRouter()
		r.HandleFunc("/event", deliveryTest.CreateEvent).Methods("POST")
		req, err := http.NewRequest("POST", "/event", bytes.NewBuffer(bodyEventJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", test.userId)
		r.ServeHTTP(w, req.WithContext(userIdContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var updateEventTests = []struct {
	id         int
	userId     string
	eventId    string
	vars       map[string]string
	event      *models.Event
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "100",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		"100",
		nil,
		nil,
		nil,
		response.ErrorResponse(response.ErrJSONDecoding.Error())},
	{3,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "100",
		},
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestUpdateEvent(t *testing.T) {
	for _, test := range updateEventTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		eventModel := new(models.Event)
		if test.event != nil {
			eventModel.ID = test.event.ID
		}

		useCaseMock.On("UpdateEvent", eventModel, test.userId).Return(test.useCaseErr)

		bodyEventJSON, err := json.Marshal(test.event)
		require.NoError(t, err, logTestMessage+"err =", err)
		if test.event == nil {
			bodyEventJSON = nil
		}

		r := mux.NewRouter()
		r.HandleFunc("/event", deliveryTest.UpdateEvent).Methods("UPDATE")
		req, err := http.NewRequest("UPDATE", "/event", bytes.NewBuffer(bodyEventJSON))
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", test.userId)
		varsContext := context.WithValue(userIdContext, "vars", test.vars)
		r.ServeHTTP(w, req.WithContext(varsContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var deleteEventTests = []struct {
	id         int
	userId     string
	eventId    string
	vars       map[string]string
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		nil,
		response.OkResponse()},
	{2,
		"1",
		"",
		nil,
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestDeleteEvent(t *testing.T) {
	for _, test := range deleteEventTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("DeleteEvent", test.eventId, test.userId).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/event", deliveryTest.DeleteEvent).Methods("DELETE")
		req, err := http.NewRequest("DELETE", "/event", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		userIdContext := context.WithValue(context.Background(), "userId", test.userId)
		varsContext := context.WithValue(userIdContext, "vars", test.vars)
		r.ServeHTTP(w, req.WithContext(varsContext))

		wTest := httptest.NewRecorder()
		response.SendResponse(wTest, test.output)
		expected := wTest.Body
		actual := w.Body
		require.Equal(t, expected, actual, logTestMessage+" "+strconv.Itoa(test.id)+" "+"error")
	}
}

var getEventByIdTests = []struct {
	id         int
	eventId    string
	vars       map[string]string
	event      *models.Event
	useCaseErr error
	output     *response.Response
}{
	{1,
		"1",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "1",
		},
		nil,
		response.EventResponse(&models.Event{
			ID: "1",
		})},
	{1,
		"1",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "1",
		},
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestGetEventById(t *testing.T) {
	for _, test := range getEventByIdTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetEventById", test.eventId).Return(test.event, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/event/{id:[0-9]+}", deliveryTest.GetEventById).Methods("GET")
		req, err := http.NewRequest("GET", "/event/"+test.eventId, nil)
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

var getEventsTests = []struct {
	id         int
	vars       map[string]string
	query      string
	eventList  []*models.Event
	useCaseErr error
	output     *response.Response
}{
	{1,
		map[string]string{
			"query":    "testQuery",
			"category": "testCategory",
			"tags":     "testTags|testTags|testTags",
		},
		"?query=testQuery&category=testCategory&tags=testTags|testTags|testTags",
		nil,
		nil,
		response.EventsListResponse(nil)},
	{2,
		map[string]string{
			"query":    "testQuery",
			"category": "testCategory",
			"tags":     "testTags|testTags|testTags",
		},
		"?query=testQuery&category=testCategory&tags=testTags|testTags|testTags",
		nil,
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		title := test.vars["query"]
		category := test.vars["category"]
		tag := test.vars["tags"]
		tags := strings.Split(tag, "|")

		useCaseMock.On("GetEvents", title, category, tags).Return(test.eventList, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/events", deliveryTest.GetEvents).
			Queries("query", "{query}", "category", "{category}", "tags", "{tags}").
			Methods("GET")
		req, err := http.NewRequest("GET", "/events"+test.query, nil)
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

var getEventsFromAuthorTests = []struct {
	id         int
	vars       map[string]string
	query      string
	eventList  []*models.Event
	useCaseErr error
	output     *response.Response
}{
	{1,
		map[string]string{
			"authorid": "123",
		},
		"?authorid=123",
		nil,
		nil,
		response.EventsListResponse(nil)},
	{2,
		map[string]string{
			"authorid": "123",
		},
		"?authorid=123",
		nil,
		error2.ErrEmptyData,
		response.ErrorResponse(error2.ErrEmptyData.Error())},
}

func TestGetEventsFromAuthor(t *testing.T) {
	for _, test := range getEventsFromAuthorTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		authorId := test.vars["authorid"]

		useCaseMock.On("GetEventsFromAuthor", authorId).Return(test.eventList, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/events", deliveryTest.GetEventsFromAuthor).
			Queries("authorid", "{authorid:[0-9]+}").
			Methods("GET")
		req, err := http.NewRequest("GET", "/events"+test.query, nil)
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
