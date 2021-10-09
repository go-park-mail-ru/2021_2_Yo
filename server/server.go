package server

import (
	authDelivery "backend/auth/delivery/http"
	authRepository "backend/auth/repository/localstorage"
	authUseCase "backend/auth/usecase"
	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	eventRepository "backend/event/repository/localstorage"
	eventUseCase "backend/event/usecase"
	"bufio"
	"net/http"
	"os"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/swaggo/http-swagger"
)

func preflight(w http.ResponseWriter, r *http.Request) {
	log.Info("In preflight")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
}

type App struct {
	authManager  *authDelivery.Delivery
	eventManager *eventDelivery.Delivery
}

func NewApp() *App {
	f, err := os.Open("auth/secret.txt")
	if err != nil {
		log.Fatal("Main : can't open file with secret keyword!", err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	secret := scanner.Text()
	if err := f.Close(); err != nil {
		log.Fatal("Main : can't close file with secret keyword!", err)
	}

	authR := authRepository.NewRepository()
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC)

	eventR := eventRepository.NewRepository()
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		authManager:  authD,
		eventManager: eventD,
	}
}

func (app *App) Run() {
	/*
	port := os.Getenv("PORT")
	if port == "" {
		log.Error("Main : PORT must be set")
		port = "8080"
	}*/

	if err := initConfig(); err != nil {
		log.Fatalf("Ошибка при инициализации конфигов, %s", err.Error())
	}

	port := viper.GetString("port")
	r := mux.NewRouter()
	
	r.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	r.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List)
	r.Methods("OPTIONS").HandlerFunc(preflight)
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//TODO: Проверить, нужно ли это?
	r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))

	log.Info("Deploying. Port: ", port)

	errServer := http.ListenAndServe(":"+port, r)
	if errServer != nil {
		log.Error("Main : ListenAndServe error: ", errServer)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}