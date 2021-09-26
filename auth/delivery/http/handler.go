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
	Status int    `json:"status"`
	Msg    string `json:"message,omitempty"`
	Name   string `json:"name"`
}

func newResponse(status int, msg string) *response {
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
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type userDataForSignin struct {
	Mail     string `json:"email"`
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

func getUserFromJSONLogin(r *http.Request) (*userDataForSignin, error) {
	userInput := new(userDataForSignin)
	//Пытаемся декодировать JSON-запрос в структуру
	//Валидность данных проверяется на фронтенде (верно?...)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

func (h *HandlerAuth) Cors(w http.ResponseWriter, r *http.Request) {
	log.Println()
	w.Write([]byte("smth"))

}

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUserInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signup_json"}`, 500)
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}
	err = h.useCase.SignUp(newUserInput.Name, newUserInput.Surname, newUserInput.Mail, newUserInput.Password)
	switch err {
	case auth.ErrUserNotFound:
		http.Error(w, `{"error":"signup_signup"}`, 500)
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
		//Возможно, будут другие случаи
	default:
		m := response{200, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
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
	username, err := h.useCase.ParseKsenia(jwtToken)
	if err != nil {
		log.Info(err)
	}
	w.WriteHeader(http.StatusOK)
	m := response{200, "smth", username}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Write(b)
	//w.Write([]byte(jwtToken))
}

func (h *HandlerAuth) Test(w http.ResponseWriter, r *http.Request) {
	log.Println("In test")
	defer r.Body.Close()
	smth := "smth"
	w.Write([]byte(smth))
}

func (h *HandlerAuth) User(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintln(w, "Главная страница")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("No cookie")
		return
	}
	log.Println("Nice cookie")
	username, err1 := h.useCase.ParseKsenia(cookie.Value)
	log.Println("After Parse")
	log.Println(username)

	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("wrong parse")
		return
	}
	w.WriteHeader(http.StatusOK)
	m := response{200, "smth", username}
	b, err := json.Marshal(m)
	if err != nil {
		log.Info(err)
	}
	w.Write(b)
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
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}
	log.Println("Nice cookie")
	username, err1 := h.useCase.Parse(cookie.Value)
	log.Println("After Parse")
	log.Println(username)

	if err1 != nil {
		m := response{404, "smth", ""}
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		w.Write(b)
		return
	}

	m := response{200, "smth", ""}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (h *HandlerAuth) MiddleWare(handler http.Handler) http.Handler {
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
