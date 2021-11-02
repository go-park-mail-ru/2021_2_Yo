package register

import (
	"backend/server"
	//sessionMiddleware "backend/session/middleware"
	"github.com/gorilla/mux"
)

func RegisterHTTPEndpoints(subRouter *mux.Router, app *server.App) {
	subRouter.HandleFunc("", app.UserManager.GetUser).Methods("GET")
	subRouter.HandleFunc("/info", app.UserManager.UpdateUserInfo).Methods("POST")
	subRouter.HandleFunc("/password", app.UserManager.UpdateUserPassword).Methods("POST")
	subRouter.HandleFunc("/avatar", app.UserManager.UpdateUserPhoto).Methods("POST")
}
