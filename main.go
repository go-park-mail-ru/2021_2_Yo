package main

import (
	deliveryAuth "backend/auth/delivery/http"
	localStorageAuth "backend/auth/repository/localstorage"
	useCaseAuth "backend/auth/usecase"
	"github.com/rs/cors"
	"net/http"
	"os"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Hello, World!")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()

	repo := localStorageAuth.NewRepositoryUserLocalStorage()
	useCase := useCaseAuth.NewUseCaseAuth(repo)
	handler := deliveryAuth.NewHandlerAuth(useCase)

	r.HandleFunc("/", handler.MainPage).Methods("GET")
	r.HandleFunc("/signup", handler.SignUp).Methods("POST")
	r.HandleFunc("/signin", handler.SignIn).Methods("POST")
	r.HandleFunc("/test", handler.Test).Methods("GET")
	r.HandleFunc("/auth", handler.Auth).Methods("GET")
	r.HandleFunc("/list", handler.List).Methods("GET")
	//Нужен метод для SignIn с методом GET

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://bmstusa.herokuapp.com"},
	})

	mainHandler := c.Handler(r)

	log.Info("Deploying. Port: ", port)
	err := http.ListenAndServe(":"+port, mainHandler)
	if err != nil {
		log.Error("main error: ", err)
	}

}
