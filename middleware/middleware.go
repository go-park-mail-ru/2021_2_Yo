package middleware

import (
	log "backend/logger"
	"backend/response"
	"backend/response/utils"
	"backend/service/session"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const logMessage = "middleware:"

type Middlewares struct {
	manager session.Manager
}

func NewMiddlewares(sm session.Manager) *Middlewares {
	return &Middlewares{
		manager: sm,
	}
}

func (m *Middlewares) Recovery(next http.Handler) http.Handler {
	message := logMessage + "Recovery:"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(message+"err =", err)
				response.SendResponse(w, response.ErrorResponse("Internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://bmstusssa.herokuapp.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Expose-Headers",
			"Accept,Accept-Encoding,X-CSRF-Token,Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(r.Method, r.RequestURI, time.Since(start))
	})
}

func (m *Middlewares) GetVars(next http.Handler) http.Handler {
	message := logMessage + "GetVars:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Debug(message+"vars =", vars)
		if vars != nil {
			varsCtx := context.WithValue(r.Context(), "vars", vars)
			next.ServeHTTP(w, r.WithContext(varsCtx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) Auth(next http.Handler) http.Handler {
	message := logMessage + "Auth:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
			return
		}
		userId, err := m.manager.Check(cookie.Value)
		if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
