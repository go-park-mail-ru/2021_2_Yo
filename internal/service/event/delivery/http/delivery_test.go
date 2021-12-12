package http

import (
	"backend/internal/models"
	"backend/internal/service/event/usecase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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
}{
	{
		1,
		"1",
		"100",
		&models.Event{
			ID: "1",
		},
		nil,
	},
	{
		2,
		"1",
		"100",
		nil,
		nil,
	},
	{
		3,
		"1",
		"100",
		&models.Event{
			ID: "1",
		},
		errors.New("test_err"),
	},
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
	}
}

var updateEventTests = []struct {
	id         int
	userId     string
	eventId    string
	vars       map[string]string
	event      *models.Event
	useCaseErr error
}{
	{
		1,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "100",
		},
		nil,
	},
	{
		2,
		"1",
		"100",
		nil,
		nil,
		nil,
	},
	{
		3,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "100",
		},
		errors.New("test_err"),
	},
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
	}
}

var deleteEventTests = []struct {
	id         int
	userId     string
	eventId    string
	vars       map[string]string
	useCaseErr error
}{
	{
		1,
		"1",
		"100",
		map[string]string{
			"id": "100",
		},
		nil,
	},
	{
		2,
		"1",
		"",
		nil,
		errors.New("test_err"),
	},
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
	}
}

var getEventByIdTests = []struct {
	id         int
	eventId    string
	vars       map[string]string
	event      *models.Event
	useCaseErr error
}{
	{
		1,
		"1",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "1",
		},
		nil,
	},
	{
		1,
		"1",
		map[string]string{
			"id": "100",
		},
		&models.Event{
			ID: "1",
		},
		errors.New("test_err"),
	},
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
	}
}

var getEventsTests = []struct {
	id         int
	vars       map[string]string
	query      string
	eventList  []*models.Event
	useCaseErr error
}{
	{
		1,
		map[string]string{
			"query":    "testQuery",
			"category": "testCategory",
			"city":     "testCity",
			"date":     "testDate",
			"tags":     "testTags|testTags|testTags",
		},
		"?query=testQuery&category=testCategory&city=testCity&date=testDate&tags=testTags|testTags|testTags",
		nil,
		nil,
	},
	{
		2,
		map[string]string{
			"query":    "testQuery",
			"category": "testCategory",
			"tags":     "testTags|testTags|testTags",
		},
		"?query=testQuery&category=testCategory&tags=testTags|testTags|testTags",
		nil,
		errors.New("test_err"),
	},
}

func TestGetEvents(t *testing.T) {
	for _, test := range getEventsTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		title := test.vars["query"]
		category := test.vars["category"]
		city := test.vars["city"]
		date := test.vars["date"]
		tag := test.vars["tags"]
		tags := strings.Split(tag, "|")

		useCaseMock.On("GetEvents", title, category, city, date, tags).Return(test.eventList, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/events", deliveryTest.GetEvents).Methods("GET")
		req, err := http.NewRequest("GET", "/events"+test.query, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getEventsFromAuthorTests = []struct {
	id         int
	vars       map[string]string
	eventList  []*models.Event
	useCaseErr error
}{
	{
		1,
		map[string]string{
			"id": "123",
		},
		nil,
		nil,
	},
	{
		2,
		map[string]string{
			"id": "123",
		},
		nil,
		errors.New("test_err"),
	},
}

func TestGetEventsFromAuthor(t *testing.T) {
	for _, test := range getEventsFromAuthorTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		authorId := test.vars["authorid"]

		useCaseMock.On("GetCreatedEvents", authorId).Return(test.eventList, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("{id:[0-9]+}", deliveryTest.GetCreatedEvents).
			Methods("GET")
		req, err := http.NewRequest("GET", authorId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getVisitedEventsTests = []struct {
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
	{
		3,
		"1",
		errors.New("test_err"),
	},
}

func TestGetVisitedEvents(t *testing.T) {
	for _, test := range getVisitedEventsTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetVisitedEvents", test.userId).Return([]*models.Event{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/{id:[0-9]+}", deliveryTest.GetVisitedEvents).Methods("GET")
		req, err := http.NewRequest("GET", "/"+test.userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getCreatedEventsTests = []struct {
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
	{
		3,
		"1",
		errors.New("test_err"),
	},
}

func TestGetCreatedEvents(t *testing.T) {
	for _, test := range getVisitedEventsTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetCreatedEvents", test.userId).Return([]*models.Event{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/{id:[0-9]+}", deliveryTest.GetCreatedEvents).Methods("GET")
		req, err := http.NewRequest("GET", "/"+test.userId, nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var visitTests = []struct {
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
		errors.New(""),
		"2",
		errors.New("test_err"),
	},
	{
		3,
		map[string]string{
			"id": "123",
		},
		errors.New(""),
		errors.New("test_err"),
	},
}

func TestVisit(t *testing.T) {
	for _, test := range visitTests {
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

		useCaseMock.On("Visit", eId, uId).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.Visit).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var unvisitTests = []struct {
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
		errors.New(""),
		"2",
		errors.New("test_err"),
	},
	{
		3,
		map[string]string{
			"id": "123",
		},
		errors.New(""),
		errors.New("test_err"),
	},
}

func TestUnvisit(t *testing.T) {
	for _, test := range unvisitTests {
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

		useCaseMock.On("Unvisit", eId, uId).Return(test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.Unvisit).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var isVisitedTests = []struct {
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
		errors.New(""),
		"2",
		errors.New("test_err"),
	},
	{
		3,
		map[string]string{
			"id": "123",
		},
		errors.New(""),
		errors.New("test_err"),
	},
}

func TestIsVisited(t *testing.T) {
	for _, test := range unvisitTests {
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

		useCaseMock.On("IsVisited", eId, uId).Return(true, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.IsVisited).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		ctxVars := context.WithValue(context.Background(), "vars", test.vars)
		ctxUserId := context.WithValue(ctxVars, "userId", test.userId)
		req = req.WithContext(ctxUserId)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

var getCitiesTests = []struct {
	id         int
	useCaseErr error
}{
	{
		1,
		nil,
	},
	{
		2,
		errors.New("test_err"),
	},
}

func TestGetCities(t *testing.T) {
	for _, test := range unvisitTests {
		useCaseMock := new(usecase.UseCaseMock)
		deliveryTest := NewDelivery(useCaseMock)

		useCaseMock.On("GetCities").Return([]string{}, test.useCaseErr)

		r := mux.NewRouter()
		r.HandleFunc("/test", deliveryTest.GetCities).Methods("GET")
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err, logTestMessage+"NewRequest error")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
