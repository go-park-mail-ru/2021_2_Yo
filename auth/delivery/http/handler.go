package http

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	"backend/response"
	"backend/response/utils"
	"backend/session"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
	"time"
)

const logMessage = "auth:delivery:http:handler:"

type Delivery struct {
	useCase        auth.UseCase
	sessionManager session.Manager
}

func NewDelivery(useCase auth.UseCase, manager session.Manager) *Delivery {
	return &Delivery{
		useCase:        useCase,
		sessionManager: manager,
	}
}

func (h *Delivery) setJwtToken(w http.ResponseWriter, jwtToken string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	cs := w.Header().Get("Set-Cookie")
	cs += "; SameSite=None"
	w.Header().Set("Set-Cookie", cs)
}

func getUserFromRequest(r *http.Request) (*models.User, error) {
	userInput := new(response.ResponseBodyUser)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	result := &models.User{
		Name:     userInput.Name,
		Surname:  userInput.Surname,
		Mail:     userInput.Mail,
		Password: userInput.Password,
		About:    userInput.About,
	}
	_, err = govalidator.ValidateStruct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//TODO!!!!
//TODO: Работа с sessionManager:
//TODO:		Ставить куки с session_id, создавать session_id и т.п.

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
	userFromRequest, err := getUserFromRequest(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	err = h.useCase.SignUp(userFromRequest)
	if !utils.CheckIfNoError(&w, err, message, http.StatusConflict) {
		return
	}
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
		return
	}
	h.setJwtToken(w, jwtToken)
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
	userFromRequest, err := getUserFromRequest(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
		return
	}
	log.Debug(message+"jwtToken =", jwtToken)
	h.setJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "User:"
	log.Debug(message + "started")
	cookie, err := r.Cookie("session_id")
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	expiredJwtToken, err := h.useCase.Logout(cookie.Value)
	h.setJwtToken(w, expiredJwtToken)
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
func (h *Delivery) User(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "User:"
	log.Debug(message + "started")
	cookie, err := r.Cookie("session_id")
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	foundUser, err := h.useCase.ParseToken(cookie.Value)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug(message+"foundUser =", foundUser)
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) GetCSRF(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetCSRF:"
	log.Debug(message + "started")
	var err error
	cookie, _ := r.Cookie("session_id")
	CSRFToken, err := h.useCase.GetCSRFToken(cookie.Value, time.Now().Add(24 * time.Hour).Unix())
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	w.Header().Set("X-CSRF-Token", CSRFToken)
}