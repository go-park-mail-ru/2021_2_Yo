package http

import (
	"backend/auth"
	"backend/models"
	//"backend/models"
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

type userDataResponse struct {
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Mail     string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

//TODO: Нормальный response для каждых случаев, когда нужно посылать ответ
type response struct {
	Status int    `json:"status"`
	Msg    string `json:"message,omitempty"`
	Name   string `json:"name"`
}

type responseError struct {
	Error string `json:"error"`
}

func makeUserDataForResponse(user *models.User) *userDataResponse {
	return &userDataResponse{
		Name:     user.Name,
		Surname:  user.Surname,
		Mail:     user.Mail,
		Password: user.Password,
	}
}

func getUserFromJSON(r *http.Request) (*userDataResponse, error) {
	userInput := new(userDataResponse)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func sendResponse(w http.ResponseWriter, responseToSend *response) {
	//TODO: Выяснить, нужно ли делать так: w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(responseToSend)
	if err != nil {
		
		log.Error("SignUp : Response error")
		
		return
	}
	w.Write(b)
}

func sendError(w http.ResponseWriter, error string) {
	//TODO: Выяснить, нужно ли делать так: w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(&responseError{Error: error})
	w.Write(b)
}

func setCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
	cs := w.Header().Get("Set-Cookie")
	cs += "; SameSite=None"
	w.Header().Set("Set-Cookie", cs)
}

func (h *HandlerAuth) setCookieWithJwtToken(w http.ResponseWriter, userMail, userPassword string) {
	
	log.Info("setCookieWithJwtToken : started")
	
	jwtToken, err := h.useCase.SignIn(userMail, userPassword)
	if err == auth.ErrUserNotFound {
		
		log.Error("SignIn : setCookieWithJwtToken error", err)
		
		sendError(w, "User not found")
		return
	}
	
	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)
	
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
	}
	setCookie(w, cookie)
	
	log.Info("setCookieWithJwtToken : ended")
	
}

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Info("SignUp : started")

	userFromRequest, err := getUserFromJSON(r)

	if err != nil {
		log.Error("SignUp : didn't get user from JSON", err)
		sendError(w, "")
		return
	}

	log.Info("SignUp : userFromRequest = ", userFromRequest)
	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname, userFromRequest.Mail, userFromRequest.Password)
	if err != nil {
		log.Error("SignUp : SignUp error", err)
		sendError(w, "User already exists")
		return
	}

	h.setCookieWithJwtToken(w, userFromRequest.Mail, userFromRequest.Password)
	sendResponse(w, &response{
		Status: http.StatusOK,
		Msg:    "Sign Up success",
		Name:   "",
	})
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
	h.setCookieWithJwtToken(w, userFromRequest.Mail, userFromRequest.Password)
	sendResponse(w, &response{
		Status: http.StatusOK,
		Msg:    "Cookie sent!",
		Name:   "",
	})
	
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
		
		sendError(w, "Error with getting cookie")
		return
	}
	
	if cookie != nil {
		log.Info("User : cookie.value = ", cookie.Value)
	}
	
	//TODO: Отладить этот момент, мб если cookie пустая, то при инициализации cookie вылезет ошибка и вызовется предыдущий if
	if cookie == nil {
		
		log.Error("User : cookie is nil")
		
		sendError(w, "No cookie sent or wrong cookie format")
		return
	}
	userID, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		
		log.Info("User : parse error", err)
		
		sendError(w, "Error with parsing token")
		return
	}
	
	log.Info("User : userID = ", userID)
	
	foundUser, err := h.useCase.GetUserById(userID)
	if err == auth.ErrUserNotFound {
		
		log.Info("User : GetUser error", err)
		
		sendError(w, "User not found")
		return
	}
	userData := makeUserDataForResponse(foundUser)
	sendResponse(w, &response{
		Status: http.StatusOK,
		Msg:    "sending name",
		Name:   userData.Name,
	})
	
	log.Info("User : ended")
	
}
