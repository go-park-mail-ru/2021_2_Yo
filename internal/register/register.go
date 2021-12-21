package register

import (
	"backend/internal/middleware"
	authHttp "backend/internal/service/auth/delivery/http"
	eventHttp "backend/internal/service/event/delivery/http"
	userHttp "backend/internal/service/user/delivery/http"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthHTTPEndpoints(r *mux.Router, delivery *authHttp.Delivery, middlewares *middleware.Middlewares) {
	r.HandleFunc("/signup", delivery.SignUp).Methods("POST")
	r.HandleFunc("/login", delivery.SignIn).Methods("POST")
	logoutHandlerFunc := http.HandlerFunc(delivery.Logout)
	r.Handle("/logout", middlewares.Auth(logoutHandlerFunc))
}

func UserHTTPEndpoints(r *mux.Router, uDelivery *userHttp.Delivery, eDelivery *eventHttp.Delivery, mws *middleware.Middlewares) {
	r.HandleFunc("/{id:[0-9]+}", uDelivery.GetUserById).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/events/favourite", eDelivery.GetVisitedEvents).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/events/created", eDelivery.GetCreatedEvents).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/subscribers", uDelivery.GetSubscribers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/subscriptions", uDelivery.GetSubscribes).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/friends", uDelivery.GetFriends).Methods("GET")

	getUserHandlerFunc := mws.Auth(http.HandlerFunc(uDelivery.GetUser))
	r.Handle("", getUserHandlerFunc).Methods("GET")

	updateUserInfoHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.UpdateUserInfo)))
	r.Handle("/info", updateUserInfoHandlerFunc).Methods("POST")

	updateUserPasswordHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.UpdateUserPassword)))
	r.Handle("/password", updateUserPasswordHandlerFunc).Methods("POST")

	subscribeHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.Subscribe)))
	r.Handle("/{id:[0-9]+}/subscription", subscribeHandlerFunc).Methods("POST")

	unsubscribeHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.Unsubscribe)))
	r.Handle("/{id:[0-9]+}/subscription", unsubscribeHandlerFunc).Methods("DELETE")

	isSubscribedHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.IsSubscribed)))
	r.Handle("/{id:[0-9]+}/subscription", isSubscribedHandlerFunc).Methods("GET")

	getAllNotificationsHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.GetAllNotifications)))
	r.Handle("/notifications/all", getAllNotificationsHandlerFunc).Methods("GET")

	getNewNotificationsHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.GetNewNotifications)))
	r.Handle("/notifications/new", getNewNotificationsHandlerFunc).Methods("GET")

	updateNotificationsStatusHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.UpdateNotificationsStatus)))
	r.Handle("/notifications/all", updateNotificationsStatusHandlerFunc).Methods("POST")

	getFriendsHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.GetFriends)))
	r.Handle("/friends", getFriendsHandlerFunc).Methods("GET")

	inviteHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(uDelivery.Invite)))
	r.Handle("/{id:[0-9]+}/invite", inviteHandlerFunc).Methods("POST")
}

func EventHTTPEndpoints(r *mux.Router, delivery *eventHttp.Delivery, mws *middleware.Middlewares) {
	//TODO: Попросить фронт заменить "query" на "title", ибо понятно, почему.
	r.HandleFunc("", delivery.GetEvents).Methods("GET")
	r.HandleFunc("/cities", delivery.GetCities).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", delivery.GetEventById).Methods("GET")
	updateEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.UpdateEvent)))
	r.Handle("/{id:[0-9]+}", updateEventHandlerFunc).Methods("POST")
	deleteEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.DeleteEvent)))
	r.Handle("/{id:[0-9]+}", deleteEventHandlerFunc).Methods("DELETE")
	createEventHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.CreateEvent)))
	r.Handle("", createEventHandlerFunc).Methods("POST")

	visitHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.Visit)))
	r.Handle("/{id:[0-9]+}/favourite", visitHandlerFunc).Methods("POST")

	unvisitHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.Unvisit)))
	r.Handle("/{id:[0-9]+}/favourite", unvisitHandlerFunc).Methods("DELETE")

	isVisitedHandlerFunc := mws.Auth(mws.GetVars(http.HandlerFunc(delivery.IsVisited)))
	r.Handle("/{id:[0-9]+}/favourite", isVisitedHandlerFunc).Methods("GET")
}
