package http

import (
	log "backend/logger"
	"backend/response"
	"backend/service/user"
	"backend/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

//TODO: Проверять везде контекст на пустоту

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	foundUser, err := h.useCase.GetUserById(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug(message+"imgUrl =", foundUser.ImgUrl)
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
	err := r.ParseMultipartForm(0)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userReader := strings.NewReader(r.FormValue("json"))
	u, err := response.GetUserFromRequest(userReader)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	imgUrl, err := utils.SaveImageFromRequest(r, "file")
	/*
		if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
		}
	*/
	u.ImgUrl = imgUrl
	u.ID = r.Context().Value("userId").(string)
	err = h.useCase.UpdateUserInfo(u)
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
	log.Debug(message + "ended")
}
