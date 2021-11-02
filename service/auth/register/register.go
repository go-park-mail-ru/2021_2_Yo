package register

import (
	"backend/server"
	sessionMiddleware "backend/session/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPEndpoints(subRouter *mux.Router, app *server.App, middleware *sessionMiddleware.Middleware) {
	subRouter.HandleFunc("/signup", app.AuthManager.SignUp).Methods("POST")
	subRouter.HandleFunc("/login", app.AuthManager.SignIn).Methods("POST")
	subRouter.Handle("/logout", middleware.Auth(http.HandlerFunc(app.AuthManager.Logout))).Methods("GET")
}
