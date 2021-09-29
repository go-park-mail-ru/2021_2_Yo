package http

import (
	"backend/auth"
	"backend/response"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerAuth struct {
	useCase auth.UseCaseAuth
}

func NewHandlerAuth(useCase auth.UseCaseAuth) *HandlerAuth {
	return &HandlerAuth{
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

func (h *HandlerAuth) setCookieWithJwtToken(w http.ResponseWriter, jwtToken string) {
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

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Info("SignUp : started")
	userFromRequest, err := getUserFromJSON(r)
	if err != nil {
		log.Error("SignUp : didn't get user from JSON", err)
		response.SendResponse(w, response.ErrorResponse("SignUp : didn't get user from JSON"))
		return
	}
	log.Info("SignUp : userFromRequest = ", userFromRequest)
	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname, userFromRequest.Mail, userFromRequest.Password)
	if err != nil {
		log.Error("SignUp : SignUp error", err)
		response.SendResponse(w, response.ErrorResponse("User already exists"))
		return
	}
	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error("SignIn : useCase.SignIn error", err)
		response.SendResponse(w, response.ErrorResponse("User not found"))
		return
	}
	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Info("SignUp : ended")
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
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
		response.SendResponse(w, response.ErrorResponse("User not found"))
		return
	}
	log.Info("SignIn : jwtToken = ", jwtToken)
	h.setCookieWithJwtToken(w, jwtToken)
	response.SendResponse(w, response.OkResponse())
	log.Info("SignIn : ended")
}

func (h *HandlerAuth) User(w http.ResponseWriter, r *http.Request) {
	log.Info("User : started")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Error("User : cookie error", err)
		response.SendResponse(w, response.ErrorResponse("Error with getting cookie"))
		return
	}
	if cookie != nil {
		log.Info("User : cookie.value = ", cookie.Value)
	}
	foundUser, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		log.Error("User : User token parsing error", err)
		response.SendResponse(w, response.ErrorResponse("Error with parsing token"))
		return
	}
	log.Info("User : Found User = ", foundUser)
	response.SendResponse(w, response.UsernameResponse(foundUser.Name))
	log.Info("User : ended")
}
