package http

import (
	"backend/auth"
	"encoding/json"
	"errors"
	"net/http"
)

type HandlerAuth struct {
	useCase auth.UseCase
}

//auth.UseCase - это чистый интерфейс
//Передаём интерфейс, а не конкретную реализацию, поскольку нужно будет передавать мок для тестирования
func NewHandlerAuth(useCase auth.UseCase) *HandlerAuth {
	return &HandlerAuth{
		useCase: useCase,
	}
}

//Структура, в которую мы попытаемся перевести JSON-запрос
type userData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getUserFromJSON(r *http.Request) (*userData, error) {
	userInput := new(userData)
	//Пытаемся декодировать JSON-запрос в структуру
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	//Если не получилось нормально распарсить
	if userInput.Username == "" || userInput.Password == "" {
		return nil, errors.New("Invalid data for new user")
	}
	return userInput, nil
}

func (h *HandlerAuth) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"signup method should be POST"}`, 400)
		return
	}
	newUserInput, err := getUserFromJSON(r)
	if err != nil {
		http.Error(w, `{"error":"signup_json"}`, 500)
		return
	}
	err = h.useCase.SignUp(newUserInput.Username, newUserInput.Password)
	if err != nil {
		http.Error(w, `{"error":"signup_signup"}`, 500)
		return
	}
}

func (h *HandlerAuth) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodPost:
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
		//По факту, я должен буду ответить на w, что всё хорошо и отправить токен.
		w.Write([]byte(jwtToken))
	case http.MethodGet:
		w.Write([]byte("Got method GET"))
	default:
		http.Error(w, `{"error":"SignIn method should be POST or GET"}`, 400)
	}
	return
}
