package middleware

import (
	"backend/response/utils"
	"backend/service/session"
	"context"
	"net/http"
)

type Middleware struct {
	manager session.Manager
}

func NewMiddleware(manager session.Manager) *Middleware {
	return &Middleware{
		manager: manager,
	}
}

const logMessage = "session:middleware:"

func (m *Middleware) Auth(next http.Handler) http.Handler {
	message := logMessage + "Auth:"
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
