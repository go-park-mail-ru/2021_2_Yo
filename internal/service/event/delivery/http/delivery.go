package http

import (
	"backend/internal/response"
	"backend/internal/service/event"
	"backend/internal/utils"
	log "backend/pkg/logger"
	"backend/pkg/notificator"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const logMessage = "service:event:delivery:http:"

type Delivery struct {
	useCase     event.UseCase
	notificator notificator.NotificationManager
}

func NewDelivery(useCase event.UseCase, notificator notificator.NotificationManager) *Delivery {
	return &Delivery{
		useCase:     useCase,
		notificator: notificator,
	}
}

func (h *Delivery) CreateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "CreateEvent:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
	err := r.ParseMultipartForm(5 << 20)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	eventReader := strings.NewReader(r.FormValue("json"))
	eventFromRequest, err := response.GetEventFromRequest(eventReader)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}

	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		response.CheckIfNoError(&w, err, message)
		return
	}
	if err == nil {
		eventFromRequest.ImgUrl = imgUrl
	}
	eventFromRequest.AuthorId = userId
	eventID, err := h.useCase.CreateEvent(eventFromRequest)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.EventIdResponse(eventID))
	_ = h.notificator.NewEventNotification(userId, eventID)
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateEvent:"
	log.Debug(message + "started")
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	eventId := vars["id"]
	userId := r.Context().Value(response.CtxString("userId")).(string)
	err := r.ParseMultipartForm(5 << 20)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	eventReader := strings.NewReader(r.FormValue("json"))
	eventFromRequest, err := response.GetEventFromRequest(eventReader)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		response.CheckIfNoError(&w, err, message)
		return
	}
	if err == nil {
		eventFromRequest.ImgUrl = imgUrl
	}
	eventFromRequest.ID = eventId
	err = h.useCase.UpdateEvent(eventFromRequest, userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "DeleteEvent:"
	log.Debug(message + "started")
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	eventId := vars["id"]
	userId := r.Context().Value(response.CtxString("userId")).(string)
	err := h.useCase.DeleteEvent(eventId, userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) GetEventById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvent:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	resultEvent, err := h.useCase.GetEventById(eventId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.EventResponse(resultEvent))
}

func (h *Delivery) GetEvents(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvents:"
	log.Debug(message + "started")
	q := r.URL.Query()

	var userId string
	var title string
	var category string
	var city string
	var tag string
	var date string

	if len(q["userId"]) > 0 {
		userId = q["userId"][0]
	}
	if len(q["query"]) > 0 {
		title = q["query"][0]
	}
	if len(q["category"]) > 0 {
		category = q["category"][0]
	}
	if len(q["tags"]) > 0 {
		tag = q["tags"][0]
	}
	if len(q["city"]) > 0 {
		city = q["city"][0]
	}
	if len(q["date"]) > 0 {
		date = q["date"][0]
	}
	tags := strings.Split(tag, "|")

	eventsList, err := h.useCase.GetEvents(userId, title, category, city, date, tags)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventsList))
}

func (h *Delivery) GetVisitedEvents(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetVisitedEvents:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	eventList, err := h.useCase.GetVisitedEvents(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventList))
}

func (h *Delivery) GetCreatedEvents(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetCreatedEvents:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	eventList, err := h.useCase.GetCreatedEvents(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventList))
}

func (h *Delivery) Visit(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Visit:"
	log.Debug(message + "started")
	vars, ok := r.Context().Value(response.CtxString("vars")).(map[string]string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	userId, ok := r.Context().Value(response.CtxString("userId")).(string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	eventId := vars["id"]
	err := h.useCase.Visit(eventId, userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Unvisit(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Unvisit:"
	log.Debug(message + "started")
	vars, ok := r.Context().Value(response.CtxString("vars")).(map[string]string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	userId, ok := r.Context().Value(response.CtxString("userId")).(string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	eventId := vars["id"]
	err := h.useCase.Unvisit(eventId, userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) IsVisited(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "IsVisited:"
	log.Debug(message + "started")
	vars, ok := r.Context().Value(response.CtxString("vars")).(map[string]string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	userId, ok := r.Context().Value(response.CtxString("userId")).(string)
	if !ok {
		response.CheckIfNoError(&w, errors.New("type casting error"), message)
	}
	eventId := vars["id"]
	res, err := h.useCase.IsVisited(eventId, userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.FavouriteResponse(res))
	log.Debug(message + "ended")
}

func (h *Delivery) GetCities(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetCities:"
	log.Debug(message + "started")
	res, err := h.useCase.GetCities()
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.CitiesResponse(res))
	log.Debug(message + "ended")
}
