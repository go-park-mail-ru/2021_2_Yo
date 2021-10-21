package middleware

import (
	log "backend/logger"
	"backend/models"
	"backend/response"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const logMessage = "auth:delivery:http:middleware:"

type Middleware struct {
	//Потом может пригодиться
	//TODO: засунуть сюда репозиторий для Auth, чтобы пользователей доставать
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

/*
СТРУКТУРА MIDDLEWARE
func (m *Middleware) MiddlewareName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//DO stuff
		next.ServeHTTP(w, r)
	})
}
*/

func (m *Middleware) Recovery(next http.Handler) http.Handler {
	message := logMessage + "Recovery:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(message+"err =", err)
				response.SendResponse(w, response.ErrorResponse("Internal server error"))
				//TODO: Разобраться, нужно ли здесь отсылать 500 через w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	message := logMessage + "Logging:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(r.Method, r.RequestURI, time.Since(start))
	})
}

func getUserFromJSON(r *http.Request) (*models.User, error) {
	userInput := new(response.ResponseBodyUser)
	err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, err
	}
	result := &models.User{
		Name:     userInput.Name,
		Surname:  userInput.Surname,
		Mail:     userInput.Mail,
		Password: userInput.Password,
	}
	return result, nil
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	message := logMessage + "Auth:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userFromRequest, err := getUserFromJSON(r)
		if err != nil {
			log.Error(message+"err =", err)
			response.SendResponse(w, response.ErrorResponse("Не получилось получить пользователя из JSON"))
			return
		}
		log.Debug(message+"user from request =", userFromRequest)
		userCtx := context.WithValue(context.Background(), "user", userFromRequest)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
