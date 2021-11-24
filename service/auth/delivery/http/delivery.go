package http

import (
	log "backend/pkg/logger"
	"backend/pkg/response"
	"backend/pkg/utils"
	"backend/service/auth"
	error2 "backend/service/auth/error"
	"backend/service/email"
	"net/http"
)

const logMessage = "service:auth:delivery:http:"

type Delivery struct {
	//С большой, так как нужен будет в других функциях
	UseCase auth.UseCase
}

func NewDelivery(useCase auth.UseCase) *Delivery {
	return &Delivery{
		UseCase: useCase,
	}
}

func setSessionIdCookie(w http.ResponseWriter, sessionId string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func setExpiredCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_id",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	log.Debug("setExpiredCooke:cookie.value =", cookie.Value)
	http.SetCookie(w, cookie)
}

func (h *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignUp:"
	log.Debug(message + "started")
	u, err := response.GetUserFromRequest(r.Body)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userId, err := h.UseCase.SignUp(u)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	sessionId, err := h.UseCase.CreateSession(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	CSRFToken, err := h.UseCase.CreateToken(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	setSessionIdCookie(w, sessionId)
	log.Info(CSRFToken)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response.SendResponse(w, response.OkResponse())
	email.SendEmail("Подтвержение регистрации", "Вы успешно зарегистрировались на BMSTUSA!", []string{u.Mail})
	log.Debug(message + "ended")
}

func (h *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignIn:"
	log.Debug(message + "started")
	u, err := response.GetUserFromRequest(r.Body)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userId, err := h.UseCase.SignIn(u)
	if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
		return
	}
	sessionId, err := h.UseCase.CreateSession(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	CSRFToken, err := h.UseCase.CreateToken(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	setSessionIdCookie(w, sessionId)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Logout:"
	log.Debug(message + "started")
	setExpiredCookie(w)
	cookie, err := r.Cookie("session_id")
	if err != nil {
		err = error2.ErrCookie
	}
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	err = h.UseCase.DeleteSession(cookie.Value)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}
