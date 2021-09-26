package http

import (
	"backend/auth"
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

type response struct {
	Status string `json:"status"`
	Msg    string `json:"message,omitempty"`
}

func newResponse(status, msg string) *response {
	return &response{
		Status: status,
		Msg:    msg,
	}
}

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
type userDataForSignup struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Mail string `json:"email"`
	Password string `json:"password"`
}

type userDataForSignin struct {
	Mail string `json:"email"`
	Password string `json:"password"`
}



func getUserFromJSON(r *http.Request) (*userDataForSignup, error) {
	userInput := new(userDataForSignup)
	//Пытаемся декодировать JSON-запрос в структуру
	//Валидность данных проверяется на фронтенде (верно?...)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func getUserFromJSONLogin (r *http.Request) (*userDataForSignin, error) {
	userInput := new(userDataForSignin)
	//Пытаемся декодировать JSON-запрос в структуру
	//Валидность данных проверяется на фронтенде (верно?...)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func(h *HandlerAuth) Cors (w http.ResponseWriter, r *http.Request){
	log.Println()
	w.Write([]byte("smth"))

}

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUserInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signup_json"}`, 500)
		return
	}
	err = h.useCase.SignUp(newUserInput.Name, newUserInput.Surname, newUserInput.Mail, newUserInput.Password)
	switch err {
	case auth.ErrUserNotFound:
		http.Error(w, `{"error":"signup_signup"}`, 500)
		return
		//Возможно, будут другие случаи
	default:
		return
	}
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userInput, err := getUserFromJSONLogin(r)
	if err != nil {
		http.Error(w, `{"error":"signin_json"}`, 500)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jwtToken, err := h.useCase.SignIn(userInput.Mail, userInput.Password)
	if err != nil {
		http.Error(w, `{"error":"signin_signin"}`, 500)
		w.WriteHeader(http.StatusNotFound)
		return
	}

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtToken))
}

func (h *HandlerAuth) Test(w http.ResponseWriter, r *http.Request) {
	log.Println("In test")
	defer r.Body.Close()
	smth := "smth"
	w.Write([]byte(smth))
}

func (h *HandlerAuth) Auth(w http.ResponseWriter, r *http.Request) {
	log.Println("In auth")
	defer r.Body.Close()
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

/*

func (h *HandlerAuth) mySignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	defer r.Body.Close()
	newUserInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signup_json"}`, 500)
		return
	}
}
*/

func (h *HandlerAuth) List(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	usernames := h.useCase.List()
	for _, username := range usernames {
		fmt.Println(username)
	}
}

func (h *HandlerAuth) MainPage(w http.ResponseWriter, r *http.Request) {
	
	defer r.Body.Close()
	fmt.Fprintln(w, "Главная страница")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No cookie")
		return
		//w.WriteHeader(http.Statu)
	}
	log.Println("Nice cookie")
	username, err1 := h.useCase.Parse(cookie.Value)
	log.Println("After Parse")
	log.Println(username)

	if err1 != nil {
		log.Println("Parse error")
		log.Println(err1)
		return
	}
	log.Println("hello " + username)
	w.Write([]byte("hello " + username))

}

func (h *HandlerAuth) MiddleWare(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Info("in middleware")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type","application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
		handler.ServeHTTP(w,r)
	})
}