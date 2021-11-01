package middleware

import (
	"backend/csrf"
	"backend/response/utils"
	"context"
	"net/http"
)

type Middleware struct {
	manager csrf.Manager
}

func NewMiddleware(manager csrf.Manager) *Middleware {
	return &Middleware{
		manager: manager,
	}
}

const logMessage = "session:middleware:"

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	message := logMessage + "CSRF:"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gottenToken := w.Header().Get("X-CSRF-Token")
		userId, err := m.manager.Check(gottenToken)
		if !utils.CheckIfNoError(&w, err, message, http.StatusNotFound) {
			return
		}
		userCtx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(userCtx))
	})
}
