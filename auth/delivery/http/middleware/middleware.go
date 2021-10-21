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
	log.Debug("Auth:Middleware:Recovery")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error("Auth:http:middleware:recovery panic error = ", err)
				response.SendResponse(w, response.ErrorResponse("Internal server error"))
				//TODO: Разобраться, нужно ли здесь отсылать 500 через w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	log.Debug("Middleware:Logging")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(r.Method, " ", r.RequestURI, " ", time.Since(start))
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	log.Debug("Middleware:Cors")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
		next.ServeHTTP(w, r)
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message := "Middleware:Auth:"
		log.Debug(message + "started")
		userFromRequest, err := getUserFromJSON(r)
		if err != nil {
			log.Error(message+"err = ", err)
			response.SendResponse(w, response.ErrorResponse("Не получилось получить пользователя из JSON"))
			return
		}
		log.Debug(message+"user from request = ", userFromRequest)
		userCtx := context.WithValue(context.Background(), "user", userFromRequest)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
