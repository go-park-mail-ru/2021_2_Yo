package server

import (
	"backend/auth"
	authDelivery "backend/auth/delivery/http"
	"backend/auth/delivery/http/middleware"
	authLocalRepository "backend/auth/repository/localstorage"
	authRepository "backend/auth/repository/postgres"
	authUseCase "backend/auth/usecase"
	_ "backend/docs"
	"backend/event"
	eventDelivery "backend/event/delivery/http"
	eventLocalRepository "backend/event/repository/localstorage"
	eventRepository "backend/event/repository/postgres"
	eventUseCase "backend/event/usecase"
	log "backend/logger"
	"bufio"
	"errors"

	//"errors"
	"fmt"
	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

const logMessage = "server:"

type App struct {
	authManager  *authDelivery.Delivery
	eventManager *eventDelivery.Delivery
	db           *sql.DB
}

func getSecret(isRemoteServer bool, pathToSecretFile string) (string, error) {
	message := logMessage + "getSecret:"
	log.Debug(message + "started")
	if isRemoteServer {
		secret := os.Getenv("SECRET")
		if secret == "" {
			secret = "secret1234"
			err := errors.New("Can't get secret from environment")
			log.Error(message+"err =", err)
			//return "", err
		}
		return secret, nil
	} else {
		f, err := os.Open(pathToSecretFile)
		if err != nil {
			log.Error(message+"err =", err)
			return "", err
		}
		scanner := bufio.NewScanner(f)
		scanner.Scan()
		secret := scanner.Text()
		if err := f.Close(); err != nil {
			log.Error(message+"err =", err)
			return "", err
		}
		return secret, nil
	}
}

func initDB(isRemoteServer bool) (*sql.DB, error) {
	message := logMessage + "initDB:"
	log.Debug(message + "started")
	if !isRemoteServer {
		return nil, nil
	}

	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	log.Debug(message+"connStr =", connStr)

	db, err := sql.Connect("postgres", connStr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	log.Info("db status =", db.Stats())
	return db, nil
}

func NewApp(isRemoteServer bool, logLevel logrus.Level) (*App, error) {
	message := logMessage + "NewApp:"
	log.Init(logLevel)
	log.Info(fmt.Sprintf(message+"started, isRemoteServer = %t, log level = %s", isRemoteServer, logLevel))
	secret, err := getSecret(isRemoteServer, "auth/secret.txt")

	db, err := initDB(isRemoteServer)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	var authR auth.Repository
	var eventR event.Repository
	if isRemoteServer {
		authR = authRepository.NewRepository(db)
		eventR = eventRepository.NewRepository(db)
	} else {
		authR = authLocalRepository.NewRepository()
		eventR = eventLocalRepository.NewRepository()
	}

	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC)

	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		authManager:  authD,
		eventManager: eventD,
		db:           db,
	}, nil
}

func Preflight(w http.ResponseWriter, r *http.Request) {
	message := logMessage + "Preflight:"
	log.Info(message + "start")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS,HEAD")
	log.Info(message + "end")
}

func newRouterWithEndpoints(app *App) *mux.Router {
	midwar := middleware.NewMiddleware()

	authMux := mux.NewRouter()
	authMux.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	authMux.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	authMux.Use(midwar.Auth)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(Preflight)
	r.Handle("/signup", authMux)
	r.Handle("/login", authMux)
	r.HandleFunc("/logout", app.authManager.Logout).Methods("GET")
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.GetEvent).Methods("GET")
	//TODO: Проверка на пользователя, отправляющего запрос
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", app.eventManager.CreateEvent).Methods("POST")
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//Сначала будет вызываться recovery, потом cors, а потом logging
	r.Use(midwar.Logging)
	r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))
	r.Use(midwar.Recovery)
	return r
}

func (app *App) Run() error {
	if app.db != nil {
		defer app.db.Close()
	}
	message := logMessage + "Run:"
	log.Info(message + "start")
	r := newRouterWithEndpoints(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(message+"port =", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	return nil
}
