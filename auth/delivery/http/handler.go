package http

import (
	"backend/auth"
	"backend/models"
	//"backend/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var cookies = make(map[string]string)

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type HandlerAuth struct {
	useCase auth.UseCase
}

func NewHandlerAuth(useCase auth.UseCase) *HandlerAuth {
	//auth.UseCase - это чистый интерфейс
	//Передаём интерфейс, а не конкретную реализацию, поскольку нужно будет передавать мок для тестирования
	return &HandlerAuth{
		useCase: useCase,
	}
}

//Структура, в которую мы попытаемся перевести JSON-запрос
//Эта структура - неполная, она, например, не содержит ID и чего-нибудь ещё (дату рождения, например)
type userDataForSignUp struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type userDataForSignIn struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type userDataForResponse struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Mail    string `json:"mail"`
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

func makeUserDataForResponse(user *models.User) *userDataForResponse {
	return &userDataForResponse{
		Name:    user.Name,
		Surname: user.Surname,
		Mail:    user.Mail,
	}
}

func getUserFromJSONSignUp(r *http.Request) (*userDataForSignUp, error) {
	userInput := new(userDataForSignUp)
	//Пытаемся декодировать JSON-запрос в структуру
	//Валидность данных проверяется на фронтенде (верно?...)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func getUserFromJSONSignIn(r *http.Request) (*userDataForSignIn, error) {
	userInput := new(userDataForSignIn)
	//Пытаемся декодировать JSON-запрос в структуру
	//Валидность данных проверяется на фронтенде (верно?...)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

//КОСТЫЛЬ
func sendRespose(w http.ResponseWriter, responseToSend *response) {
	////////////////////////////////
	b, err := json.Marshal(responseToSend)
	if err != nil {
		/////////
		log.Error("SignUp : Response error")
		/////////
		return
	}
	w.Write(b)
	/////////////////////////////////
}

//Не уверен, что здесь указатель, проверить!
func (h *HandlerAuth) setCookieWithJwtToken(w http.ResponseWriter, userMail, userPassword string) {
	/////////
	log.Info("setCookieWithJwtToken : started")
	/////////
	//TODO: Сделать так, чтобы SignIn возвращал только токен и ошибку. информация о user будет возвращаться в User (функция)
	//TODO: Вроде сделал
	jwtToken, err := h.useCase.SignIn(userMail, userPassword)
	/////////
	log.Info("setCookieWithJwtToken : jwtToken = ", jwtToken)
	/////////
	if err == auth.ErrUserNotFound {
		/////////
		log.Error("SignIn : setCookieWithJwtToken error")
		/////////
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
	}
	//Костыль, добавляем ещё одну куку, которая не записывается голангом
	http.SetCookie(w, cookie)
	cs := w.Header().Get("Set-Cookie")
	cs += "; SameSite=None"
	w.Header().Set("Set-Cookie", cs)
	/////////
	log.Info("setCookieWithJwtToken : ended")
	/////////
}

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	/////////
	log.Info("SignUp : started")
	/////////
	userFromRequest, err := getUserFromJSONSignUp(r)
	/////////
	log.Info("SignUp : userFromRequest = ", userFromRequest)
	/////////
	if err != nil {
		/////////
		log.Error("SignUp : didn't get user from JSON")
		/////////
		return
	}
	err = h.useCase.SignUp(userFromRequest.Name, userFromRequest.Surname, userFromRequest.Mail, userFromRequest.Password)
	if err != nil {
		/////////
		log.Error("SignUp : SignUp error")
		/////////
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(&responseError{Error: "User already exists"})
		w.Write(b)
		return
	}
	//TODO: Поставить Cookie с jwt-токеном при регистрации
	//TODO: Вроде сделал
	h.setCookieWithJwtToken(w, userFromRequest.Mail, userFromRequest.Password)
	w.WriteHeader(http.StatusOK)
	sendRespose(w, &response{200, "smth", ""})
	/////////
	log.Info("SignUp : ended")
	/////////
	return
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
	/////////
	log.Info("SignIn : started")
	/////////
	defer r.Body.Close()
	userFromRequest, err := getUserFromJSONSignIn(r)
	/////////
	log.Info("SignIn : userFromRequest = ", userFromRequest)
	/////////
	if err != nil {
		/////////
		log.Error("SignIn : getUserFromJSON error")
		/////////
		return
	}
	_, err1 := h.useCase.GetUser(userFromRequest.Mail, userFromRequest.Password)
	if err1 == auth.ErrUserNotFound {
		/////////
		log.Error("SignIn : GetUser error")
		/////////
		/////////////////////////////////////////////////
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(&responseError{Error: "User not found"})
		w.Write(b)
		/////////////////////////////////////////////////
		//sendRespose(w, &response{301, "User not found!", ""})
		return
	}
	h.setCookieWithJwtToken(w, userFromRequest.Mail, userFromRequest.Password)
	w.WriteHeader(http.StatusOK)
	sendRespose(w, &response{200, "Cookie sent!", ""})
	////////////////////////////////
	/*
	m := response{200, "smth", ""}
	b, err := json.Marshal(m)
	if err != nil {
		/////////
		log.Error("SignIn : Response error")
		/////////
		return
	}
	w.Write(b)
	 */
	/////////////////////////////////
	/////////
	log.Info("SignIn : ended")
	/////////
	return
}

func (h *HandlerAuth) MiddleWare(handler http.Handler) http.Handler {
	/////////
	log.Info("MiddleWare : started & ended")
	/////////
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
	/////////
	log.Info("User : started")
	/////////
	//TODO: Отправить информацию о пользователе таким образом:
	/*
		//Получаем данные о пользователе для того, чтобы отправить их пользователю
		userData := makeUserDataForResponse(foundUser)
		w.WriteHeader(http.StatusOK)
		userDataToWrite, err := json.Marshal(userData)
		if err != nil {
			/////////
			log.Error("User : json.Marshall error")
			/////////
			return
		}
		w.Write(userDataToWrite)
	 */

	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	/////////
	if cookie != nil {
		log.Info("User : cookie.value = ", cookie.Value)
	}
	/////////
	if err != nil {
		/////////
		log.Error("User : getting cookie error")
		/////////
		w.WriteHeader(http.StatusTeapot)
		return
	}
	//TODO: Разобраться, как работает ParseToken и что возвращает
	userID, err := h.useCase.ParseToken(cookie.Value)
	/////////
	log.Info("User : userID = ", userID)
	/////////
	if err != nil {
		/////////
		log.Info("User : parse error")
		/////////
		w.WriteHeader(http.StatusTeapot)
		return
	}
	//TODO: отправить информацию пользователю
	foundUser, err := h.useCase.GetUserById(userID)
	if err == auth.ErrUserNotFound {
		/////////
		log.Info("User : GetUser error")
		/////////
		w.WriteHeader(http.StatusTeapot)
		return
	}
	///////////////////
	userData := makeUserDataForResponse(foundUser)
	sendRespose(w, &response{
		Status: http.StatusOK,
		Msg:    "sending name",
		Name:   userData.Name,
	})
	/*
	w.WriteHeader(http.StatusOK)
	userDataToWrite, err := json.Marshal(userData)
	if err != nil {
		/////////
		log.Error("User : json.Marshall error")
		/////////
		return
	}
		sendRespose(w, &response{
			Status: http.StatusOK,
			Msg:    string(userDataToWrite),
			Name:   "",
		})
	 */

	/////////
	log.Info("User : ended")
	/////////
}

//TODO: 1. User: response(status, Name, Surname, Mail)
//TODO: 2. Login: если неправильно, отправлять response(status, body : {error : "Didn't login"})