package http

import (
	"backend/auth"
	"backend/models"
	"fmt"
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

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	/////////
	log.Debug("SignUp : started")
	/////////
	newUserInput, err := getUserFromJSONSignUp(r)
	if err != nil {
		/////////
		log.Error("SignUp : didn't get user from JSON")
		/////////
		http.Error(w, `{"error":"signup_json"}`, 500)
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}
	fmt.Println(newUserInput)
	err = h.useCase.SignUp(newUserInput.Name, newUserInput.Surname, newUserInput.Mail, newUserInput.Password)
	if err != nil {
		log.Error("SignUp : SignUp error")
		http.Error(w, `{"error":"signup_signup"}`, 500)
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}
	/////////
	log.Debug("SignUp : ended")
	/////////
	return
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
	/////////
	log.Debug("SignIn : started")
	/////////
	defer r.Body.Close()
	userInput, err := getUserFromJSONSignIn(r)
	if err != nil {
		/////////
		log.Error("SignIn : getUserFromJSON error")
		/////////
		http.Error(w, `{"error":"signin_json"}`, 500)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	foundUser, jwtToken, err := h.useCase.SignIn(userInput.Mail, userInput.Password)
	if err == auth.ErrUserNotFound {
		/////////
		log.Error("SignIn : useCase.SignIn error")
		/////////
		http.Error(w, `{"error":"signin_user_not_found"}`, 500)
		w.WriteHeader(http.StatusNotFound)
		return
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
	//Получаем данные о пользователе для того, чтобы отправить их пользователю
	userData := makeUserDataForResponse(foundUser)
	w.WriteHeader(http.StatusOK)

	userDataToWrite, err := json.Marshal(userData)
	if err != nil {
		/////////
		log.Error("SignIn : json.Marshall error")
		/////////
		return
	}
	w.Write(userDataToWrite)
	/////////
	log.Debug("SignIn : ended")
	/////////
	return
}

func (h *HandlerAuth) Auth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("In auth")
	kukan, _ := r.Cookie("auth")
	log.Println(cookies)
	log.Println(kukan.Value)
	if kukan.Value == "" {
		log.Println("in error")
		cookies["rarara"] = "Blabla"
		//w.WriteHeader(http.StatusNotFound)
		cookie := http.Cookie{Name: "auth", Value: "rarara", Secure: true}
		http.SetCookie(w, &cookie)
		cs := w.Header().Get("Set-Cookie")
		cs += "; SameSite=None"
		w.Header().Set("Set-Cookie", cs)
		log.Println(w.Header().Get("Set-Cookie"))
		log.Println(cookie.Value)
	} else {
		log.Println(kukan.Value)
		_, ok := cookies[kukan.Value]
		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			//w.WriteHeader(http.StatusBadGateway)
		}
	}
}

func (h *HandlerAuth) List(w http.ResponseWriter, r *http.Request) {
	/////////
	log.Debug("List : started")
	/////////
	fmt.Println("")
	fmt.Println("=============================")
	fmt.Println("=========U==S==E==R==S=======")
	fmt.Println("=============================")
	defer r.Body.Close()
	users := h.useCase.List()
	for _, user := range users {
		fmt.Println(user)
		userData := makeUserDataForResponse(&user)
		userDataToWrite, _ := json.Marshal(userData)
		w.Write(userDataToWrite)
	}
	fmt.Println("=============================")
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	fmt.Println("=============================")
	fmt.Println("")
	/////////
	log.Debug("List : ended")
	/////////
	return
}

func (h *HandlerAuth) MainPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintln(w, "Главная страница")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No cookie")
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}

	username, err := h.useCase.ParseToken(cookie.Value)
	if err != nil {
		log.Println("Parse error", err)
		return
	}
	log.Println("hello " + username)
	w.Write([]byte("hello " + username))
}

func (h *HandlerAuth) MiddleWare(handler http.Handler) http.Handler {
	/////////
	log.Debug("MiddleWare : started & ended")
	/////////

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("in middleware")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
		handler.ServeHTTP(w, r)
	})
}
