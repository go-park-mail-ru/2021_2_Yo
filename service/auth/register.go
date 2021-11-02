package auth

import (
	"backend/service/auth/delivery/http"
	"github.com/gorilla/mux"
)

func RegisterHTTPEndpoints(r *mux.Router, delivery *http.Delivery) {
	r.HandleFunc("/signup", delivery.SignUp)
}
