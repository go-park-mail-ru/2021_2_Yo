package server

import (
	authDelivery "backend/auth/delivery/http"
	authRepository "backend/auth/repository/postgres"
	authUseCase "backend/auth/usecase"
	eventRepository "backend/event/repository/postgres"
	eventUseCase "backend/event/usecase"
	"bufio"
	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"

	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type App struct {
	authManager  *authDelivery.Delivery
	eventManager *eventDelivery.Delivery
	db *sql.DB
}

func preflight(w http.ResponseWriter, r *http.Request) {
	log.Info("In preflight")
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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("main : Can't open DB", err)
		return nil, err
	}
	log.Println("db status = ", db.Stats())
	return db, nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func NewApp() (*App, error) {
	secret := getSecret("auth/secret.txt")
	user := viper.GetString("db.username")
	password := viper.GetString("db.password")

	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")
	
	connStr := "user="+user+" password="+password+" dbname="+dbname+" sslmode="+sslmode

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

	return &App{
		authManager:  authD,
		eventManager: eventD,
		db: db,
	}, nil
}

func (app *App) Run() error {
	if err := initConfig(); err != nil {
		log.Error("Server : Run() initConfig err", err)
	}
	defer app.db.Close()

	port := viper.GetString("port")
	r := mux.NewRouter()

	r.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	r.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List)
	r.Methods("OPTIONS").HandlerFunc(preflight)
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//TODO: Проверить, нужно ли это? Или preflight достаточно?
	r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))

	log.Info("Server : Run() : Deploying, port = ", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("Server : Run() : ListenAndServe error: ", err)
		return err
	}
	return nil
}
