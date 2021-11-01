package http

import (
	"backend/event"
	log "backend/logger"
	"backend/response"
	"backend/response/utils"
	"github.com/gorilla/mux"
	"net/http"
)

const logMessage = "event:delivery:http:"

type Delivery struct {
	useCase event.UseCase
}

func NewDelivery(useCase event.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

//@Summmary List
//@Tags Events
//@Description "Список мероприятий"
//@Produce json
//@Success 200 {object} response.ResponseBodyEventList
//@Failure 404 {object} response.BaseResponse
//@Router /events [get]
func (h *Delivery) List(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "List:"
	log.Debug(message + "started")
	eventsList, err := h.useCase.List()
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.EventsListResponse(eventsList))
}

func (h *Delivery) GetEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvent:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	resultEvent, err := h.useCase.GetEvent(eventId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug("delivery:getEvent:resultEvent.authorId = ", resultEvent.AuthorId)
	response.SendResponse(w, response.EventResponse(resultEvent))
}

func (h *Delivery) CreateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "CreateEvent:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	eventFromRequest, err := response.GetEventFromJSON(r)
	log.Debug(message+"eventFromRequest = ", *eventFromRequest)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	eventID, err := h.useCase.CreateEvent(eventFromRequest, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.EventIdResponse(eventID))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateEvent:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	eventId := vars["id"]
	userId := r.Context().Value("userId").(string)
	eventFromRequest, err := response.GetEventFromJSON(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	eventFromRequest.ID = eventId
	err = h.useCase.UpdateEvent(eventFromRequest, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "DeleteEvent:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	eventId := vars["id"]
	userId := r.Context().Value("userId").(string)
	err := h.useCase.DeleteEvent(eventId, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}
