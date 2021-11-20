package http

import (
	log "backend/pkg/logger"
	"backend/response"
	"backend/service/event"
	"backend/utils"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const logMessage = "service:event:delivery:http:"

type Delivery struct {
	useCase event.UseCase
}

func NewDelivery(useCase event.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) CreateEvent(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "CreateEvent:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	err := r.ParseMultipartForm(5 << 20)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	eventReader := strings.NewReader(r.FormValue("json"))
	eventFromRequest, err := response.GetEventFromRequest(eventReader)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		utils.CheckIfNoError(&w, err, message, http.StatusBadRequest)
		return
	}
	if err == nil {
		eventFromRequest.ImgUrl = imgUrl
	}
	eventFromRequest.AuthorId = userId
	eventID, err := h.useCase.CreateEvent(eventFromRequest)
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
	err := r.ParseMultipartForm(5 << 20)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	eventReader := strings.NewReader(r.FormValue("json"))
	eventFromRequest, err := response.GetEventFromRequest(eventReader)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		utils.CheckIfNoError(&w, err, message, http.StatusBadRequest)
		return
	}
	if err == nil {
		eventFromRequest.ImgUrl = imgUrl
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

func (h *Delivery) GetEventById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvent:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	eventId := vars["id"]
	resultEvent, err := h.useCase.GetEventById(eventId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug("delivery:getEvent:resultEvent.authorId = ", resultEvent.AuthorId)
	response.SendResponse(w, response.EventResponse(resultEvent))
}

type getEventsVars struct {
	title    string   `valid:"type(string),length(0|50)" san:"xss"`
	category string   `valid:"type(string),length(0|50)" san:"xss"`
	tags     []string `valid:"type(string),length(0|50)" san:"xss"`
}

func (h *Delivery) GetEvents(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEvents:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	title := vars["query"]
	category := vars["category"]
	tag := vars["tags"]
	tags := strings.Split(tag, "|")
	//err := response.ValidateAndSanitize(title)
	//if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
	//	return
	//}
	eventsList, err := h.useCase.GetEvents(title, category, tags)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventsList))
}

func (h *Delivery) GetEventsFromAuthor(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEventsFromAuthor:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	authorId := vars["authorid"]
	eventList, err := h.useCase.GetEventsFromAuthor(authorId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventList))
}

/*
GET /user/14/visited
*/
func (h *Delivery) Visit(w http.ResponseWriter, r *http.Request) {

	message := logMessage + "Visit:"
	log.Debug(message + "started")

	vars, ok := r.Context().Value("vars").(map[string]string)
	if !ok {
		utils.CheckIfNoError(&w, errors.New("type casting error"), message, http.StatusInternalServerError)
	}
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		utils.CheckIfNoError(&w, errors.New("type casting error"), message, http.StatusInternalServerError)
	}
	eventId := vars["id"]

	err := h.useCase.Visit(eventId, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}

	response.SendResponse(w, response.OkResponse())

	log.Debug(message + "ended")

}

func (h *Delivery) GetVisitedEvents(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetEventsFromAuthor:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	eventList, err := h.useCase.GetVisitedEvents(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.EventListResponse(eventList))
}
