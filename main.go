package main

import (
	deliveryAuth "backend/auth/delivery/http"
	localStorageAuth "backend/auth/repository/localstorage"
	useCaseAuth "backend/auth/usecase"
	"github.com/rs/cors"
	"net/http"
	//"os"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Preflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
}

func main() {
	log.Println("Hello, World!")
	/*
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}*/

	r := mux.NewRouter()

	repo := localStorageAuth.NewRepositoryUserLocalStorage()
	useCase := useCaseAuth.NewUseCaseAuth(repo)
	handler := deliveryAuth.NewHandlerAuth(useCase)
	r.Use(handler.MiddleWare)

	r.HandleFunc("/", handler.MainPage).Methods("GET")
	r.HandleFunc("/signup", handler.Cors).Methods("OPTIONS")
	r.HandleFunc("/signup", handler.SignUp).Methods("POST")
	r.HandleFunc("/signin", handler.SignIn).Methods("POST")
	r.HandleFunc("/test", handler.Test).Methods("GET")
	r.HandleFunc("/auth", handler.Auth).Methods("GET")
	r.HandleFunc("/list", handler.List).Methods("GET")
	r.Methods("OPTIONS").HandlerFunc(Preflight)
	//Нужен метод для SignIn с методом GET

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://bmstusa.herokuapp.com"},
	})

	mainHandler := c.Handler(r)

	log.Info("Deploying. Port: ", )
	err := http.ListenAndServe(":8080", mainHandler)
	if err != nil {
		log.Error("main error: ", err)
	}

}
