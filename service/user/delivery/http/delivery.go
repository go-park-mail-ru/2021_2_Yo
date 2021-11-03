package http

import (
	log "backend/logger"
	"backend/response"
	"backend/response/utils"
	"backend/service/image"
	"backend/service/user"
	"net/http"
)

const logMessage = "service:user:delivery:http:"

type Delivery struct {
	useCase    user.UseCase
	imgManager image.Manager
}

func NewDelivery(useCase user.UseCase, imgManager image.Manager) *Delivery {
	return &Delivery{
		useCase:    useCase,
		imgManager: imgManager,
	}
}

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	foundUser, err := h.useCase.GetUser(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) GetUserById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUserById:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	foundUser, err := h.useCase.GetUser(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	u, err := response.GetUserFromRequest(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	err = h.useCase.UpdateUserInfo(userId, u.Name, u.Surname, u.About)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	u, err := response.GetUserFromRequest(r)
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

func (h *Delivery) UpdateUserPhoto(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserPhoto:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	r.ParseMultipartForm(1 << 2)
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Error("error Retrieving the File")
		return
	}
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	defer file.Close()

	log.Info("Uploaded File: %+v\n", handler.Filename)
	log.Info("File Size: %+v\n", handler.Size)
	log.Info("MIME Header: %+v\n", handler.Header)

	err = h.imgManager.SaveFile(userId, handler.Filename, file)
	if err != nil {
		log.Error(err)
	}
	log.Info(w, "Successfully Uploaded File\n"+"")
}
