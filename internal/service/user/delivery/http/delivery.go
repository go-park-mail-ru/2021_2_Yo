package http

import (
	response "backend/internal/response"
	"backend/internal/service/user"
	"backend/internal/utils"
	log "backend/pkg/logger"
	"backend/pkg/notificator"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const logMessage = "service:user:delivery:http:"

type Delivery struct {
	useCase     user.UseCase
	notificator notificator.NotificationManager
}

func NewDelivery(useCase user.UseCase, notificator notificator.NotificationManager) *Delivery {
	return &Delivery{
		useCase:     useCase,
		notificator: notificator,
	}
}

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
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
	userFromRequest.ID = r.Context().Value(response.CtxString("userId")).(string)
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
	userId := r.Context().Value(response.CtxString("userId")).(string)
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

func (h *Delivery) GetFriends(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetFriends:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
	subscribers, err := h.useCase.GetFriends(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetVisitors(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetVisitors:"
	log.Debug(message + "started")
	eventId := r.Context().Value(response.CtxString("eventId")).(string)
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
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	subscriberId := r.Context().Value(response.CtxString("userId")).(string)
	subscribedId := vars["id"]
	err := h.useCase.Subscribe(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	err = h.notificator.NewSubscriberNotification(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Unsubscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	subscriberId := r.Context().Value(response.CtxString("userId")).(string)
	subscribedId := vars["id"]
	err := h.useCase.Unsubscribe(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	_ = h.notificator.DeleteSubscribeNotification(subscribedId, subscriberId)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) IsSubscribed(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "IsSubscribed:"
	log.Debug(message + "started")
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	subscriberId := r.Context().Value(response.CtxString("userId")).(string)
	subscribedId := vars["id"]
	res, err := h.useCase.IsSubscribed(subscribedId, subscriberId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.SubscribedResponse(res))
	log.Debug(message + "ended")
}

func (h *Delivery) Invite(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Invite:"
	log.Debug(message + "started")
	q := r.URL.Query()
	var eventId string
	if len(q["eventId"]) > 0 {
		eventId = q["eventId"][0]
	}
	vars := r.Context().Value(response.CtxString("vars")).(map[string]string)
	userId := r.Context().Value(response.CtxString("userId")).(string)
	receiverId := vars["id"]
	err := h.notificator.InvitationNotification(receiverId, userId, eventId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetAllNotifications:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
	res, err := h.notificator.GetAllNotifications(userId)
	log.Debug(res)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.NotificationListResponse(res))
	log.Debug(message + "ended")
}

func (h *Delivery) GetNewNotifications(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetNewNotifications:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
	res, err := h.notificator.GetNewNotifications(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.NotificationListResponse(res))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateNotificationsStatus(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateNotificationsStatus:"
	log.Debug(message + "started")
	userId := r.Context().Value(response.CtxString("userId")).(string)
	err := h.notificator.UpdateNotificationsStatus(userId)
	if !response.CheckIfNoError(&w, err, message) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

