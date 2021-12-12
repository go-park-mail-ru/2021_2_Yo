package http

import (
	"backend/internal/notification"
	response2 "backend/internal/response"
	"backend/internal/service/user"
	"backend/internal/utils"
	log "backend/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const logMessage = "service:user:delivery:http:"

type Delivery struct {
	useCase     user.UseCase
	notificator notification.Notificator
}

func NewDelivery(useCase user.UseCase, notificator notification.Notificator) *Delivery {
	return &Delivery{
		useCase:     useCase,
		notificator: notificator,
	}
}

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	foundUser, err := h.useCase.GetUserById(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	CSRFToken, err := utils.GenerateCsrfToken(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response2.SendResponse(w, response2.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) GetUserById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUserById:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	foundUser, err := h.useCase.GetUserById(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	err := r.ParseMultipartForm(5 << 20)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	userReader := strings.NewReader(r.FormValue("json"))
	userFromRequest, err := response2.GetUserFromRequest(userReader)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		response2.CheckIfNoError(&w, err, message)
		return
	}
	if err == nil {
		userFromRequest.ImgUrl = imgUrl
	}
	userFromRequest.ID = r.Context().Value("userId").(string)
	err = h.useCase.UpdateUserInfo(userFromRequest)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserPassword:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	u, err := response2.GetUserFromRequest(r.Body)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	err = h.useCase.UpdateUserPassword(userId, u.Password)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribers(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribes(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribes:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribes(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetVisitors(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetVisitors:"
	log.Debug(message + "started")
	eventId := r.Context().Value("eventId").(string)
	userList, err := h.useCase.GetVisitors(eventId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.UserListResponse(userList))
	log.Debug(message + "ended")
}

func (h *Delivery) Subscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Subscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	err := h.useCase.Subscribe(subscribedId, subscriberId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	subscriber, err := h.useCase.GetUserById(subscriberId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	err = h.notificator.NewSubscriber(subscribedId, subscriber.Name)
	if err != nil {
		//To db
		//storeNotification()
	}
	response2.SendResponse(w, response2.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Unsubscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	err := h.useCase.Unsubscribe(subscribedId, subscriberId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) IsSubscribed(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "IsSubscribed:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	res, err := h.useCase.IsSubscribed(subscribedId, subscriberId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.SubscribedResponse(res))
	log.Debug(message + "ended")
}
