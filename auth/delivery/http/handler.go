package http

import (
	"backend/auth"
	"backend/response"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Delivery struct {
	useCase auth.UseCase
}

func NewDelivery(useCase auth.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func getUserFromJSON(r *http.Request) (*response.ResponseBodyUser, error) {
	userInput := new(response.ResponseBodyUser)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
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
	log.Info("SignUp : started")
	userFromRequest, err := getUserFromJSON(r)
	if err != nil {
		log.Error("SignUp : didn't get user from JSON", err)
		response.SendResponse(w, response.ErrorResponse("Не получилось получить пользователя из JSON"))
		return
	}
	log.Info("SignUp : userFromRequest = ", userFromRequest)
	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname, userFromRequest.Mail, userFromRequest.Password)
	if err != nil {
		log.Error("SignUp : SignUp error", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь уже зарегестрирован"))
		return
	}
	log.Println("Auth:Handler:Signup: mail, pass ", userFromRequest.Mail, userFromRequest.Password)
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error("SignIn : useCase.SignIn error", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь не найден"))
		return
	}
	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Info("SignUp : ended")
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
	log.Info("SignIn : started")
	userFromRequest, err := getUserFromJSON(r)
	if err != nil {
		log.Error("SignIn : getUserFromJSON error", err)
		return
	}
	log.Info("SignIn : userFromRequest = ", userFromRequest)
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error("SignIn : useCase.SignIn error", err)
		response.SendResponse(w, response.ErrorResponse("Пользователь не найден"))
		return
	}
	log.Info("SignIn : jwtToken = ", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Info("SignIn : ended")
}

//@Summmary User
//@Tags auth
//@Description "Главная страница"
//@Produce json
//@Success 200 {object} response.BaseResponse
//@Failure 404 {object} response.BaseResponse
//@Router /user [get]
func (h *Delivery) User(w http.ResponseWriter, r *http.Request) {
	log.Info("User : started")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Error("User : cookie error", err)
		response.SendResponse(w, response.ErrorResponse("Ошибка с получением Cookie"))
		return
	}
	if cookie != nil {
		log.Info("User : cookie.value = ", cookie.Value)
	}
	foundUser, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		log.Error("User : User token parsing error", err)
		response.SendResponse(w, response.ErrorResponse("Ошибка с парсингом токена"))
		return
	}
	log.Info("User : Found User = ", foundUser)
	response.SendResponse(w, response.UsernameResponse(foundUser.Name))
	log.Info("User : ended")
}
