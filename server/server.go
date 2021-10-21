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

	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	log.Info(connStr)

	db, err := initDB(connStr)
	if err != nil {
		log.Error("NewApp : initDB error", err)
		return nil, err
	}

	authR := authRepository.NewRepository(db)
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC)

	eventR := eventRepository.NewRepository(db)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	log.Init(logrus.InfoLevel)

	return &App{
		authManager:  authD,
		eventManager: eventD,
		db:           db,
	}, nil
}

func (app *App) Run() error {
	defer app.db.Close()

	port := viper.GetString("port")
	r := mux.NewRouter()
	r.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	r.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List)
	r.Methods("OPTIONS").HandlerFunc(preflight)
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	midwar := middleware.NewMiddleware()
	//Сначала будет вызываться recovery, потом cors, а потом logging
	r.Use(midwar.Logging)
	r.Use(midwar.CORS)
	r.Use(midwar.Recovery)

	log.Info("Server:Run():Deploying, port = ", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("Server:Run():ListenAndServe error: ", err)
		return err
	}
	return nil
}
