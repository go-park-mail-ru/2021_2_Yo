package server

import (
	authDelivery "backend/auth/delivery/http"
	"backend/auth/delivery/http/middleware"
	authRepository "backend/auth/repository/postgres"
	authUseCase "backend/auth/usecase"
	eventRepository "backend/event/repository/postgres"
	eventUseCase "backend/event/usecase"
	"bufio"
	"fmt"
	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"

	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"net/http"

	log "backend/logger"
	"github.com/spf13/viper"
)

const logMessage = "server:"

type App struct {
	authManager  *authDelivery.Delivery
	eventManager *eventDelivery.Delivery
	db           *sql.DB
}

func getSecret(pathToSecretFile string) (string, error) {
	message := logMessage + "getSecret:"
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

func initDB() (*sql.DB, error) {
	message := logMessage + "initDB:"

	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	log.Debug(message+"connStr = ", connStr)

	db, err := sql.Connect("postgres", connStr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	log.Info("db status =", db.Stats())
	return db, nil
}

func NewApp() (*App, error) {
	message := logMessage + "NewApp:"
	log.Init(logrus.DebugLevel)
	secret, err := getSecret("auth/secret.txt")
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	db, err := initDB()
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	authR := authRepository.NewRepository(db)
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC)

	eventR := eventRepository.NewRepository(db)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		authManager:  authD,
		eventManager: eventD,
		db:           db,
	}, nil
}

func (app *App) Run() error {
	defer app.db.Close()

	message := logMessage + "Run:"
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
auth:
usecase:
usecase:
	r.Use(midwar.Recovery)

	port := viper.GetString("port")
	log.Info(message+"port =", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error(message+"err =", err)
		return err
	}
	return nil
}
