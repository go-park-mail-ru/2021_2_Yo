package main

import (
	deliveryAuth "backend/auth/delivery/http"
	localStorageAuth "backend/auth/repository/localstorage"
	useCaseAuth "backend/auth/usecase"
	"net/http"
	"os"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Preflight(w http.ResponseWriter, r *http.Request) {
	log.Info("In preflight")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
}

func main() {

	log.SetFormatter(new(log.JSONFormatter))
	log.Println("Hello, World!")

	port := os.Getenv("PORT")
	if port == "" {
		//log.Error("$PORT must be set")
		port = "8080"
	}

	r := mux.NewRouter()

	repo := localStorageAuth.NewRepositoryUserLocalStorage()
	useCase := useCaseAuth.NewUseCaseAuth(repo)
	handler := deliveryAuth.NewHandlerAuth(useCase)

	r.HandleFunc("/signup", handler.SignUp).Methods("POST")
	r.HandleFunc("/signin", handler.SignIn).Methods("POST")
	r.HandleFunc("/user", handler.User).Methods("GET")
	r.Methods("OPTIONS").HandlerFunc(Preflight)

	r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))

	log.Info("Deploying. Port: ", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("main error: ", err)
	}

}
