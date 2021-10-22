package http

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	"backend/response"
	"github.com/asaskevich/govalidator"
	"net/http"
)

const logMessage = "auth:delivery:http:handler:"

type Delivery struct {
	useCase auth.UseCase
}

func NewDelivery(useCase auth.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) setCookieWithJwtToken(w http.ResponseWriter, jwtToken string) {
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
	userFromRequest := r.Context().Value("user").(*models.User)

	_, err := govalidator.ValidateStruct(userFromRequest)
	if err != nil {
		log.Error(message+"err =", err)
	}

	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname, userFromRequest.Mail, userFromRequest.Password)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь уже зарегестрирован"))
		return
	}
	log.Debug(message+"mail, pass = ", userFromRequest.Mail, userFromRequest.Password)
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь не найден"))
		return
	}
	log.Debug(message+"jwtToken =", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
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
	userFromRequest := r.Context().Value("user").(*response.ResponseBodyUser)

	_, err := govalidator.ValidateStruct(userFromRequest)
	if err != nil {
		log.Error(message+"err =", err)
	}

	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь не найден"))
		return
	}
	log.Debug(message+"jwtToken =", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Debug(message + "ended")
}

func (h *Delivery) Logout(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Logout:"
	log.Debug(message + "started")
	userFromRequest := r.Context().Value("user").(*response.ResponseBodyUser)

	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь не найден"))
		return
	}
	jwtToken = "Null"
	log.Debug(message+"jwtToken =", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
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
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Ошибка с получением Cookie"))
		return
	}
	foundUser, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		log.Error(message+"err =", err)
		response.SendResponse(w, response.ErrorResponse("Ошибка с парсингом токена"))
		return
	}
	log.Debug(message+"foundUser =", foundUser)
	response.SendResponse(w, response.UsernameResponse(foundUser.Name))
	log.Debug(message + "ended")
}
