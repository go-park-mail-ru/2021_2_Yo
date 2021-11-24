package http

import (
	log "backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/utils"
	"backend/service/email"
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
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	CSRFToken, err := utils.GenerateCsrfToken(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
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
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	err := r.ParseMultipartForm(5 << 20)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userReader := strings.NewReader(r.FormValue("json"))
	userFromRequest, err := response.GetUserFromRequest(userReader)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	if err == utils.ErrFileExt {
		utils.CheckIfNoError(&w, err, message, http.StatusBadRequest)
		return
	}
	if err == nil {
		userFromRequest.ImgUrl = imgUrl
	}
	userFromRequest.ID = r.Context().Value("userId").(string)
	err = h.useCase.UpdateUserInfo(userFromRequest)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
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
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	err = h.useCase.UpdateUserPassword(userId, u.Password)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	email.SendEmail("Ваш пароль был изменён", "Если это были не вы, обратитесь в службу безопасности,возможно, ваш аккаунт собираются угнать", []string{u.Mail})
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribers(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetSubscribes(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribes(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserListResponse(subscribers))
	log.Debug(message + "ended")
}

func (h *Delivery) GetVisitors(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Visit:"
	log.Debug(message + "started")
	eventId := r.Context().Value("eventId").(string)
	userList, err := h.useCase.GetVisitors(eventId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
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
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Subscribe:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	subscriberId := r.Context().Value("userId").(string)
	subscribedId := vars["id"]
	err := h.useCase.Unsubscribe(subscribedId, subscriberId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
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
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	response.SendResponse(w, response.SubscribedResponse(res))
	log.Debug(message + "ended")
}
