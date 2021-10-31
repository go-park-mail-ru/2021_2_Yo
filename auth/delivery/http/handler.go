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
	"github.com/gorilla/mux"
	"net/http"
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

func (h *Delivery) setSessionIdCookie(w http.ResponseWriter, sessionId string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	cs := w.Header().Get("Set-Cookie")
	cs += "; SameSite=None"
	w.Header().Set("Set-Cookie", cs)
}

func getUserFromRequest(r *http.Request) (*models.User, error) {
	userInput := new(models.ResponseBodyUser)
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
	userId, err := h.useCase.SignUp(userFromRequest)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	sessionId, err := h.sessionManager.Create(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	h.setSessionIdCookie(w, sessionId)
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
	userFromRequest, err := getUserFromRequest(r)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	userId, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
		return
	}
	sessionId, err := h.sessionManager.Create(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusInternalServerError) {
		return
	}
	h.setSessionIdCookie(w, sessionId)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Logout:"
	log.Debug(message + "started")
	cookie, err := r.Cookie("session_id")
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	err = h.sessionManager.Delete(cookie.Value)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
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
	log.Debug(message+"foundUser =", *foundUser)
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) GetUserById(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "GetUserWithId:"
	log.Debug(message + "started")
	vars := mux.Vars(r)
	userId := vars["id"]
	foundUser, err := h.useCase.GetUser(userId)
	if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
		return
	}
	log.Debug(message+"foundUser =", *foundUser)
	response.SendResponse(w, response.UserResponse(foundUser))
	log.Debug(message + "ended")
}

func (h *Delivery) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "UpdateUserInfo:"
	log.Debug(message + "started")
	userId := r.Context().Value("userId").(string)
	u, err := getUserFromRequest(r)
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
	u, err := getUserFromRequest(r)
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
