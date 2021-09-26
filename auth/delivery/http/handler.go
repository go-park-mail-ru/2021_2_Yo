package http

import (
	"backend/auth"
	"encoding/json"
	"log"
	"net/http"
)

var cookies = make(map[string]string)

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
	Username string `json:"username"`
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

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	defer r.Body.Close()
	newUserInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signup_json"}`, 500)
		return
	}
	err = h.useCase.SignUp(newUserInput.Username, newUserInput.Password)
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
	userInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signin_json"}`, 500)
		return
	}
	jwtToken, err := h.useCase.SignIn(userInput.Username, userInput.Password)
	if err != nil {
		http.Error(w, `{"error":"signin_signin"}`, 500)
		return
	}
	w.Write([]byte(jwtToken))
}

func (h *HandlerAuth) Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin",r.Header.Get("Origin"))

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Println("In test")
	defer r.Body.Close()
	smth := "smth"
	w.Write([]byte(smth))
}

func (h *HandlerAuth) Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	log.Println("In auth")
	defer r.Body.Close()
	kukan, err := r.Cookie("auth")
	log.Println(cookies)
	//log.Println(kukan.Value)
	if err != nil {
		log.Println("in error")
		cookies["rarara"] = "Blabla"
		//w.WriteHeader(http.StatusNotFound)
		cookie := http.Cookie{Name: "auth", Value: "rarara",SameSite: http.SameSiteNoneMode , Secure: true}
		http.SetCookie(w, &cookie)
		log.Println(cookie.Value)
	} else {
		log.Println(kukan.Value)
		_, ok := cookies[kukan.Value]
		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}

	}

}
