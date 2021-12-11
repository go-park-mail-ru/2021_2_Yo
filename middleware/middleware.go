package middleware

import (
	log "backend/pkg/logger"
	"backend/pkg/response"
	"backend/service/auth"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const logMessage = "middleware:"

var allowedOrigins = []string{"", "http://127.0.0.1:3000", "https://bmstusssa.herokuapp.com"}

type Middlewares struct {
	authService auth.UseCase
}

func NewMiddlewares(authService auth.UseCase) *Middlewares {
	return &Middlewares{
		authService: authService,
	}
}

func (m *Middlewares) Recovery(next http.Handler) http.Handler {
	message := logMessage + "Recovery:"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(message+"err = ", err)
				response.SendResponse(w, response.StatusResponse(http.StatusInternalServerError))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isAllowed := false
		for _, o := range allowedOrigins {
			if origin == o {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
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
		if r.RequestURI != "/metrics" {
			log.Info(r.Method+" "+r.RequestURI+" ", time.Since(start))
		}
	})
}

func (m *Middlewares) GetVars(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if !response.CheckIfNoError(&w, err, message) {
			return
		}
		userId, err := m.authService.CheckSession(cookie.Value)
		if !response.CheckIfNoError(&w, err, message) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}

func (m *Middlewares) CSRF(next http.Handler) http.Handler {
	message := logMessage + "CSRF:"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gottenToken := (*r).Header.Get("X-CSRF-Token")
		userId, err := m.authService.CheckToken(gottenToken)
		if !response.CheckIfNoError(&w, err, message) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
