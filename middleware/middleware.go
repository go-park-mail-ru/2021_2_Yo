package middleware

import (
	log "backend/logger"
	"backend/response"
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

func (m *Middleware) CORS(next http.Handler) http.Handler {
	message := logMessage + "CORS:"
	log.Debug(message + "started")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,csrf-token,Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS,HEAD")
		//TODO: Попросить фронт не присылать options
		if r.Method == http.MethodOptions {
			return
		}
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
