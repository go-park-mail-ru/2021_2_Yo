package http

import (
	"backend/event"
	log "backend/logger"
	"backend/models"
	"backend/response"
	"encoding/json"
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

func (h *Delivery) GetEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvent:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	if eventId == "" {
		err := errors.New("eventId is null")
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get eventId out of url"))
		return
	}
	resultEvent, err := h.useCase.GetEvent(eventId)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get event with ID"))
		return
	}
	response.SendResponse(w, response.EventResponse(resultEvent))
}

func getEventFromJSON(r *http.Request) (*models.Event, error) {
	eventInput := new(response.ResponseBodyEvent)
	err := json.NewDecoder(r.Body).Decode(eventInput)
	if err != nil {
		return nil, err
	}
	result := &models.Event{
		Title:       eventInput.Title,
		Description: eventInput.Description,
		Text:        eventInput.Text,
		City:        eventInput.City,
		Category:    eventInput.Category,
		Viewed:      eventInput.Viewed,
		Tag:         eventInput.Tag,
		Date:        eventInput.Date,
		Geo:         eventInput.Geo,
	}
	return result, nil
}

func (h *Delivery) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateEvent:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	if eventId == "" {
		err := errors.New("eventId is null")
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get eventId out of url"))
		return
	}
	eventFromRequest, err := getEventFromJSON(r)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get event from JSON"))
		return
	}
	//TODO: Validate struct
	err = h.useCase.UpdateEvent(eventId, eventFromRequest)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't update such event"))
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) CreateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SetEvent:"
	log.Debug(message + "started")
	eventFromRequest, err := getEventFromJSON(r)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't get event from JSON"))
		return
	}
	//TODO: Validate struct
	eventID, err := h.useCase.CreateEvent(eventFromRequest)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Can't create such event"))
		return
	}
	response.SendResponse(w, response.EventIdResponse(eventID))
	log.Debug(message + "ended")
}
