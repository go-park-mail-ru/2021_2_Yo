package server

import (
	authDelivery "backend/auth/delivery/http"
	"backend/auth/delivery/http/middleware"
	authRepository "backend/auth/repository/localstorage"
	authUseCase "backend/auth/usecase"
	eventRepository "backend/event/repository/localstorage"
	eventUseCase "backend/event/usecase"
	"bufio"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"

	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"net/http"

	log "backend/logger"
	gorilla_handlers "github.com/gorilla/handlers"
)

type App struct {
	authManager  *authDelivery.Delivery
	eventManager *eventDelivery.Delivery
	db           *sql.DB
}

func preflight(w http.ResponseWriter, r *http.Request) {
	log.Info("Server:preflight")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS,HEAD")
}

func getSecret(pathToSecretFile string) string {
	f, err := os.Open(pathToSecretFile)
	if err != nil {
		log.Fatal("Server : can't open file with secret keyword!", err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	secret := scanner.Text()
	if err := f.Close(); err != nil {
		log.Fatal("Server : can't close file with secret keyword!", err)
	}
	return secret
}

func initDB(connStr string) (*sql.DB, error) {
	db, err := sql.Connect("postgres", connStr)
	if err != nil {
		log.Error("main : Can't open DB", err)
		return nil, err
	}
	log.Info("db status = ", db.Stats())
	return db, nil
}

func NewApp() (*App, error) {
	secret := getSecret("auth/secret.txt")

	authR := authRepository.NewRepository()
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC)

	eventR := eventRepository.NewRepository()
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	log.Init(logrus.DebugLevel)

	return &App{
		authManager:  authD,
		eventManager: eventD,
		db:           nil,
	}, nil
}

func (app *App) Run() error {
	log.Debug("Server:Run()")
	midwar := middleware.NewMiddleware()

	authMux := mux.NewRouter()
	authMux.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	authMux.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	authMux.Use(midwar.Auth)

	r := mux.NewRouter()

	r.Handle("/signup", authMux)
	r.Handle("/login", authMux)
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List).Methods("GET")
	//r.Methods("OPTIONS").HandlerFunc(preflight)
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//Сначала будет вызываться recovery, потом cors, а потом logging
	r.Use(midwar.Logging)
	//r.Use(midwar.CORS)
	r.Use(midwar.Recovery)

	r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))

	port := os.Getenv("PORT")
	if port == "" {
		log.Error("Server : PORT must be set")
		port = "8080"
	}
	log.Info("Server:Run():Deploying, port = ", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("Server:Run():ListenAndServe error: ", err)
		return err
	}
	return nil
}
