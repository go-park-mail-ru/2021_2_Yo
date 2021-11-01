package http

import (
	"backend/auth"
	error2 "backend/auth/error"
	"backend/csrf"
	log "backend/logger"
	"backend/response"
	"backend/response/utils"
	"backend/session"
	"net/http"
)

const logMessage = "auth:delivery:http:handler:"

type Delivery struct {
	useCase        auth.UseCase
	sessionManager session.Manager
	csrfManager    csrf.Manager
}

func NewDelivery(useCase auth.UseCase, manager session.Manager, csrfManager csrf.Manager) *Delivery {
	return &Delivery{
		useCase:        useCase,
		sessionManager: manager,
		csrfManager:    csrfManager,
	}
}

func setSessionIdCookie(w http.ResponseWriter, sessionId string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
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
	}
	http.SetCookie(w, cookie)
}
/*
func setCSRFCokkie(w http.ResponseWriter, csrfToken string) {
	cookie := &http.Cookie{
		Name:     "csrf-token",
		Value:    csrfToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
}
*/
//@Summmary SignUp
//@Tags auth
//@Description Регистрация
//@Accept json
//@Produce json
//@Param input body response.ResponseBodyUser true "Account Info"
//@Success 200 {object} response.Response{body=response.ResponseBodyUser}
//@Failure 404 {object} response.BaseResponse
//@Router /signup [post]
func (h *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignUp:"
	log.Debug(message + "started")
	u, err := response.GetUserFromRequest(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userId, err := h.useCase.SignUp(u)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
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
	//setCSRFCokkie(w,CSRFToken)
	log.Info(CSRFToken)
	w.Header().Set("X-CSRF-Token", CSRFToken)

	log.Debug(message+"userId =", userId)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

//@Summmary SignIn
//@Tags auth
//@Description "Авторизация"
//@Accept json
//@Produce json
//@Param input body response.ResponseBodyUser true "Account Info"
//@Success 200 {object} response.Response{body=response.ResponseBodyUser}
//@Failure 404 {object} response.BaseResponse
//@Router /signin [post]
func (h *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "SignIn:"
	log.Debug(message + "started")
	u, err := response.GetUserFromRequest(r)
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
	//setCSRFCokkie(w,CSRFToken)
	log.Info(CSRFToken)
	w.Header().Set("X-CSRF-Token", CSRFToken)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Logout:"
	log.Debug(message + "started")
	defer setExpiredCookie(w)
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
	setExpiredCookie(w)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

//@Summmary User
//@Tags auth
//@Description "Главная страница"
//@Produce json
//@Success 200 {object} response.BaseResponse
//@Failure 404 {object} response.BaseResponse
//@Router /user [get]
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
	message := logMessage + "GetUserWithId:"
	log.Debug(message + "started")
	vars := r.Context().Value("vars").(map[string]string)
	userId := vars["id"]
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
