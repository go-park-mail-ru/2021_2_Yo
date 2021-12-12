package http

import (
	"backend/internal/models"
	response2 "backend/internal/response"
	"backend/internal/service/auth"
	error2 "backend/internal/service/auth/error"
	"backend/internal/service/email"
	log "backend/pkg/logger"
	"github.com/spf13/viper"
	"net/http"
)

const logMessage = "service:auth:delivery:http:"

type Delivery struct {
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
	http.SetCookie(w, cookie)
}

func (h *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignUp:"
	log.Debug(message + "started")
	u, err := response2.GetUserFromRequest(r.Body)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	userId, err := h.UseCase.SignUp(u)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	sessionId, err := h.UseCase.CreateSession(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	CSRFToken, err := h.UseCase.CreateToken(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	template := viper.GetString("reg_html")
	setSessionIdCookie(w, sessionId)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response2.SendResponse(w, response2.OkResponse())
	info := &models.Info{
		Name: u.Name,
	}
	email.SendEmail("Подтвержение регистрации", template, []*models.Info{info})
	log.Debug(message + "ended")
}

func (h *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignIn:"
	log.Debug(message + "started")
	u, err := response2.GetUserFromRequest(r.Body)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	userId, err := h.UseCase.SignIn(u)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	sessionId, err := h.UseCase.CreateSession(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	CSRFToken, err := h.UseCase.CreateToken(userId)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	setSessionIdCookie(w, sessionId)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response2.SendResponse(w, response2.OkResponse())
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
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	err = h.UseCase.DeleteSession(cookie.Value)
	if !response2.CheckIfNoError(&w, err, message) {
		return
	}
	response2.SendResponse(w, response2.OkResponse())
	log.Debug(message + "ended")
}
