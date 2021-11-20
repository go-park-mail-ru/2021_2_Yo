package http

import (
	log "backend/logger"
	"backend/response"
	"backend/service/email"
	microAuth "backend/service/microservices/auth"
	"backend/service/user"
	"backend/utils"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const logMessage = "service:user:delivery:http:"

type Delivery struct {
	useCase     user.UseCase
	authService microAuth.AuthService
}

func NewDelivery(useCase user.UseCase, authService microAuth.AuthService) *Delivery {
	return &Delivery{
		useCase:     useCase,
		authService: authService,
	}
}

//TODO: Проверять везде контекст на пустоту

func (h *Delivery) GetUser(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUser:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	log.Debug(message+"userId =", userId)
	foundUser, err := h.useCase.GetUserById(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug(message+"imgUrl =", foundUser.ImgUrl)
	ctx := context.Background()
	CSRFToken, err := h.authService.CreateToken(ctx, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	log.Info(CSRFToken)
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
	log.Debug(message+"maxMemory =", 5<<20)
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

/*
POST /user/14/subscribe
*/
func (h *Delivery) Subscribe(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Subscribe:"
	log.Debug(message + "started")
	vars, ok := r.Context().Value("vars").(map[string]string)
	if !ok {
		utils.CheckIfNoError(&w, errors.New("type casting error"), message, http.StatusInternalServerError)
	}
	userId := r.Context().Value("userId").(string)
	subscriptedId := vars["id"]
	err := h.useCase.Subscribe(subscriptedId, userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

/*
func (h *Delivery) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetSubscribers:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	subscribers, err := h.useCase.GetSubscribers(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}
*/
