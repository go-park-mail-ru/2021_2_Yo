package middleware

import (
	log "backend/logger"
	"backend/service/jwtCsrf"
	"backend/utils"
	"context"
	"net/http"
)

type Middleware struct {
	manager jwtCsrf.Manager
}

func NewMiddleware(manager jwtCsrf.Manager) *Middleware {
	return &Middleware{
		manager: manager,
	}
}

const logMessage = "csrf:middleware:"

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	message := logMessage + "CSRF:"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gottenToken := (*r).Header.Get("X-CSRF-Token")
		log.Info("gottenToken", gottenToken)
		userId, err := m.manager.Check(gottenToken)
		if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
