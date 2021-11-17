package http

import (
	log "backend/logger"
	"backend/response"
	"backend/service/auth"
	error2 "backend/service/auth/error"
	"backend/service/csrf"
	"backend/service/email"
	microAuth "backend/service/microservices/auth"
	"backend/service/session"
	"backend/utils"
	"context"
	"net/http"
)

const logMessage = "service:auth:delivery:http:"

type Delivery struct {
	useCase        auth.UseCase
	authService    microAuth.AuthService
	sessionManager session.Manager
	csrfManager    csrf.Manager
}

func NewDelivery(useCase auth.UseCase, authService microAuth.AuthService, manager session.Manager, csrf csrf.Manager) *Delivery {
	return &Delivery{
		useCase:        useCase,
		authService:    authService, 
		sessionManager: manager,
		csrfManager:    csrf,
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
	log.Debug(message+"u =", *u)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	//userId, err := h.useCase.SignUp(u)
	ctx := context.Background()
	userResponse, err := h.authService.Create(ctx,*u)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	userId := userResponse.ID
	sessionId, err := h.sessionManager.Create(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	CSRFToken, err := h.csrfManager.Create(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	setSessionIdCookie(w, sessionId)
	log.Info(CSRFToken)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response.SendResponse(w, response.OkResponse())
	email.SendEmail("Подтвержение регистрации","Вы зарегистрировались на bmstuse",[]string{u.Mail})
	log.Debug(message + "ended")
}

func (h *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignIn:"
	log.Debug(message + "started")
	u, err := response.GetUserFromRequest(r.Body)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userId, err := h.useCase.SignIn(u.Mail, u.Password)
	if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
		return
	}
	sessionId, err := h.sessionManager.Create(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	CSRFToken, err := h.csrfManager.Create(userId)
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
	err = h.sessionManager.Delete(cookie.Value)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	CSRFToken := w.Header().Get("X-CSRF-Token")
	err = h.csrfManager.Delete(CSRFToken)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}
