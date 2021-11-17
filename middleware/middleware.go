package middleware

import (
	log "backend/logger"
	"backend/response"
	"backend/utils"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	microAuth "backend/service/microservices/auth"
)

const logMessage = "middleware:"

type Middlewares struct {
	authService    microAuth.AuthService
}

func NewMiddlewares(authService microAuth.AuthService) *Middlewares {
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
		if !utils.CheckIfNoError(&w, err, message, http.StatusBadRequest) {
			return
		}
		log.Debug(message+"cookie.value =", cookie.Value)
		ctx := context.Background()
		userId, err := m.authService.CheckSession(ctx, cookie.Value)
		if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
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
		log.Info("gottenToken", gottenToken)
		ctx := context.Background()
		userId, err := m.authService.CheckToken(ctx, gottenToken)
		if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
