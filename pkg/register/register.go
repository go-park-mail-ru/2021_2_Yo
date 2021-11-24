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

func UserHTTPEndpoints(r *mux.Router, uDelivery *userHttp.Delivery, eDelivery *eventHttp.Delivery, mws *middleware.Middlewares) {
	r.HandleFunc("/{id:[0-9]+}", uDelivery.GetUserById).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/events/favourite", eDelivery.GetVisitedEvents).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/events/created", eDelivery.GetCreatedEvents).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/subscribers", uDelivery.GetSubscribers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/subscriptions", uDelivery.GetSubscribes).Methods("GET")

	getUserHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.GetUser)))
	r.Handle("", getUserHandlerFunc).Methods("GET")

	updateUserInfoHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.UpdateUserInfo)))
	r.Handle("/info", updateUserInfoHandlerFunc).Methods("POST")

	updateUserPasswordHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.UpdateUserPassword)))
	r.Handle("/password", updateUserPasswordHandlerFunc).Methods("POST")
	//

	subscribeHandleFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.Subscribe)))
	r.Handle("/{id:[0-9]}/subscription", subscribeHandleFunc).Methods("POST")

	unsubscribeHandleFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.Unsubscribe)))
	r.Handle("/{id:[0-9]}/subscription", unsubscribeHandleFunc).Methods("DELETE")

	isSubscribedHandleFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.IsSubscribed)))
	r.Handle("/{id:[0-9]}/subscription", isSubscribedHandleFunc).Methods("GET")
	//
}

func EventHTTPEndpoints(r *mux.Router, delivery *eventHttp.Delivery, mws *middleware.Middlewares) {
	//TODO: Попросить фронт заменить "query" на "title", ибо понятно, почему.
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
	//
	visitHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.Visit)))
	r.Handle("/{id:[0-9]+}/favourite", visitHandlerFunc).Methods("POST")

	unvisitHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.Visit)))
	r.Handle("/{id:[0-9]+}/favourite", unvisitHandlerFunc).Methods("DELETE")

	isVisitedHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.Visit)))
	r.Handle("/{id:[0-9]+}/favourite", isVisitedHandlerFunc).Methods("GET")
}
