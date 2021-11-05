package register

import (
	"backend/middleware"
	authHttp "backend/service/auth/delivery/http"
	eventHttp "backend/service/event/delivery/http"
	userHttp "backend/service/user/delivery/http"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthHTTPEndpoints(r *mux.Router, delivery *authHttp.Delivery, middlewares *middleware.Middlewares) {
	r.HandleFunc("/signup", delivery.SignUp)
	r.HandleFunc("/login", delivery.SignIn)
	logoutHandlerFunc := http.HandlerFunc(delivery.Logout)
	r.Handle("/logout", middlewares.Auth(logoutHandlerFunc))
}

func useMiddlewares(r *mux.Router, path string, handlerFunc http.HandlerFunc, middlewares ...mux.MiddlewareFunc) *mux.Router {
	result := r.NewRoute().Subrouter()
	result.HandleFunc(path, handlerFunc)
	for _, mw := range middlewares {
		result.Use(mw)
	}
	return result
}

func EventHTTPEndpoints(r *mux.Router, delivery *eventHttp.Delivery, mws *middleware.Middlewares) {
	r.HandleFunc("", delivery.GetEventsFromAuthor).Queries("authorid", "{authorid:[0-9]+}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("query", "{query}", "category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("query", "{query}", "category", "{category}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("query", "{query}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("query", "{query}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("category", "{category}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Queries("tags", "{tags}").Methods("GET")
	r.HandleFunc("", delivery.GetEvents).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", delivery.GetEventById).Methods("GET")
	updateEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.UpdateEvent)))
	r.Handle("/{id:[0-9]+}", updateEventHandlerFunc).Methods("POST")
	deleteEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.DeleteEvent)))
	r.Handle("/{id:[0-9]+}", deleteEventHandlerFunc).Methods("DELETE")
	createEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.CreateEvent)))
	r.Handle("", createEventHandlerFunc).Methods("POST")
}

func UserHTTPEndpoints(r *mux.Router, delivery *userHttp.Delivery, middlewares *middleware.Middlewares) {
	r.HandleFunc("/{id:[0-9]+}", delivery.GetUserById).Methods("GET")
	r.Handle("", useMiddlewares(r, "", delivery.GetUser, middlewares.GetVars, middlewares.Auth)).Methods("POST")
	r.HandleFunc("", delivery.GetUser).Methods("GET")
	r.Handle("/info", useMiddlewares(r, "/info", delivery.UpdateUserInfo, middlewares.GetVars, middlewares.Auth)).Methods("POST")
	r.Handle("/password", useMiddlewares(r, "/password", delivery.UpdateUserPassword, middlewares.GetVars, middlewares.Auth)).Methods("POST")
	r.Handle("/avatar", useMiddlewares(r, "/avatar", delivery.UpdateUserAvatar, middlewares.GetVars, middlewares.Auth)).Methods("POST")
}
