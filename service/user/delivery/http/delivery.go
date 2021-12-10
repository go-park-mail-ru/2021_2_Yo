package http

import (
	log "backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/utils"
	"backend/service/user"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const logMessage = "service:user:delivery:http:"

type Delivery struct {
	useCase user.UseCase
}

func NewDelivery(useCase user.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	foundUser, err := h.useCase.GetUserById(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	CSRFToken, err := utils.GenerateCsrfToken(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) GetUserById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUserById:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	foundUser, err := h.useCase.GetUserById(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	err := r.ParseMultipartForm(5 << 20)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	userReader := strings.NewReader(r.FormValue("json"))
	userFromRequest, err := response.GetUserFromRequest(userReader)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		response.CheckIfNoError(&w, err, message)
		return
	}
	if err == nil {
		userFromRequest.ImgUrl = imgUrl
	}
	userFromRequest.ID = r.Context().Value("userId").(string)
	err = h.useCase.UpdateUserInfo(userFromRequest)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserPassword:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	u, err := response.GetUserFromRequest(r.Body)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	err = h.useCase.UpdateUserPassword(userId, u.Password)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribers(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribes(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribes:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribes(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetVisitors(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetVisitors:"
	log.Debug(message + "started")
	eventId := r.Context().Value("eventId").(string)
	userList, err := h.useCase.GetVisitors(eventId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.UserListResponse(userList))
	log.Debug(message + "ended")
}

func (h *Delivery) Subscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Subscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	err := h.useCase.Subscribe(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Unsubscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	err := h.useCase.Unsubscribe(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) IsSubscribed(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "IsSubscribed:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	res, err := h.useCase.IsSubscribed(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.SubscribedResponse(res))
	log.Debug(message + "ended")
}
