package http

import (
	"backend/event"
	log "backend/logger"
	"backend/response"
	"errors"
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
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get list of events"))
		return
	}
	response.SendResponse(w, response.EventsListResponse(eventsList))
}

func (h *Delivery) Event(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Event:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	if eventId == "" {
		err := errors.New("eventId is null")
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get eventId out of url"))
		return
	}
	resultEvent, err := h.useCase.Event(eventId)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get event with ID"))
		return
	}
	response.SendResponse(w, response.EventResponse(resultEvent))
}
