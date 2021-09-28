package http

import (
	"backend/auth"
	"backend/response"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerAuth struct {
	useCase auth.UseCase
}

func NewHandlerAuth(useCase auth.UseCase) *HandlerAuth {
	return &HandlerAuth{
		useCase: useCase,
	}
}


func getUserFromJSON(r *http.Request) (*response.ResponseBody, error) {
	userInput := new(response.ResponseBody)
	log.Info(r.Body)
	err := json.NewDecoder(r.Body).Decode(userInput)
	log.Info("UserInputBody = ", userInput)
	log.Info("UserInputBodyName = ",userInput.Name)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func sendResponse(w http.ResponseWriter, response *response.Response) {
	w.WriteHeader(http.StatusOK)
	log.Info("what is this ",response)
	b, err := json.Marshal(response)
	log.Info("what is this ",string(b))
	if err != nil {
		log.Error("Cound't Marshal")
		return
	}
	w.Write(b)
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
		sendResponse(w, response.ErrorResponse("SignUp : didn't get user from JSON"))
		return
	}

	log.Info("SignUp : userFromRequest = ", userFromRequest)
	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname,
		userFromRequest.Mail, userFromRequest.Password)

	if err != nil {
		log.Error("SignUp : SignUp error", err)
		sendResponse(w, response.ErrorResponse("User already exists"))
		return
	}

	jwtToken, err := h.useCase.SignIn(userFromRequest.Mail, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error("SignIn : setCookieWithJwtToken error", err)
		sendResponse(w, response.ErrorResponse("User not found"))
		return
	}

	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)

	h.setCookieWithJwtToken(w, jwtToken)

	sendResponse(w, response.OkResponse())
	log.Info("SignUp : ended")
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {

	log.Info("SignIn : started")

	userFromRequest, err := getUserFromJSON(r)

	log.Info("SignIn : userFromRequest = ", userFromRequest)

	if err != nil {

		log.Error("SignIn : getUserFromJSON error")

		return
	}

	jwtToken, err := h.useCase.SignIn(userFromRequest.Name, userFromRequest.Password)
	if err == auth.ErrUserNotFound {
		log.Error("SignIn : setCookieWithJwtToken error", err)
		sendResponse(w, response.ErrorResponse("User not found"))
		return
	}

	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)

	h.setCookieWithJwtToken(w, jwtToken)

	sendResponse(w, response.OkResponse())

	log.Info("SignIn : ended")

}

func (h *HandlerAuth) MiddleWare(handler http.Handler) http.Handler {

	log.Info("MiddleWare : started & ended")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
		handler.ServeHTTP(w, r)
	})
}

func (h *HandlerAuth) User(w http.ResponseWriter, r *http.Request) {

	log.Info("User : started")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Error("User : getting cookie error", err)
		sendResponse(w, response.ErrorResponse("Error with getting cookie"))
		return
	}

	if cookie != nil {
		log.Info("User : cookie.value = ", cookie.Value)
	}

	//TODO: Отладить этот момент, мб если cookie пустая, то при инициализации cookie вылезет ошибка и вызовется предыдущий if
	if cookie == nil {
		log.Error("User : cookie is nil")
		sendResponse(w, response.ErrorResponse("No cookie sent or wrong cookie format"))
		return
	}

	userID, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		log.Info("User : parse error", err)
		sendResponse(w, response.ErrorResponse("Error with parsing token"))
		return
	}

	log.Info("User : userID = ", userID)

	foundUser, err := h.useCase.GetUserById(userID)
	if err == auth.ErrUserNotFound {
		log.Info("User : GetUser error", err)
		sendResponse(w, response.ErrorResponse("User not found"))
		return
	}
	log.Info("User : Found User = ", foundUser)
	sendResponse(w, response.UsernameResponse(foundUser.Name))

	log.Info("User : ended")

}
