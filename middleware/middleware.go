package middleware

import (
	"backend/csrf"
	log "backend/logger"
	"backend/response"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const logMessage = "middleware:"

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Recovery(next http.Handler) http.Handler {
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

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://bmstusssa.herokuapp.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,csrf-token,Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS,HEAD")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(r.Method, r.RequestURI, time.Since(start))
	})
}

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	message := logMessage + "CSRF:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CSRFToken := w.Header().Get("X-CSRF-Token")
		log.Info(CSRFToken)
		cookie, _ := r.Cookie("session_id")
		isValidCSRFToken, _ := csrf.Token.Check(cookie.Value, CSRFToken)
		if !isValidCSRFToken {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) GetVars(next http.Handler) http.Handler {
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
