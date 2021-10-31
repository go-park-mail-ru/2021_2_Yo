package server

import (
	authDelivery "backend/auth/delivery/http"
	authRepository "backend/auth/repository/postgres"
	authUseCase "backend/auth/usecase"
	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	eventRepository "backend/event/repository/postgres"
	eventUseCase "backend/event/usecase"
	log "backend/logger"
	"backend/middleware"
	"backend/session"
	sessionMiddleware "backend/session/middleware"
	sessionRepository "backend/session/repository"
	"backend/utils"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

const logMessage = "server:"

type App struct {
	authManager    *authDelivery.Delivery
	eventManager   *eventDelivery.Delivery
	sessionManager *session.Manager
	db             *sql.DB
}

func NewApp(logLevel logrus.Level) (*App, error) {
	message := logMessage + "NewApp:"
	log.Init(logLevel)
	log.Info(fmt.Sprintf(message+"started, log level = %s", logLevel))

	secret, err := utils.GetSecret()
	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	//TODO: Параметры поменять для НЕ локал хоста
	redisAddr := flag.String("addr", "redis://user:@redis_db:6379/0", "redis addr")
	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	sessionR := sessionRepository.NewRepository(redisConn)
	sessionM := session.NewManager(*sessionR)
	authR := authRepository.NewRepository(db)
	eventR := eventRepository.NewRepository(db)
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC, *sessionM)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		authManager:    authD,
		eventManager:   eventD,
		sessionManager: sessionM,
		db:             db,
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
	mw := middleware.NewMiddleware()
	sessionMW := sessionMiddleware.NewMiddleware(*app.sessionManager)

	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	authRouter.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	authRouter.Use(sessionMW.Auth)

	r := mux.NewRouter()
	//TODO: Попросить фронт не отправлять options
	r.Methods("OPTIONS").HandlerFunc(Preflight)
	r.Handle("/signup", authRouter)
	r.Handle("/login", authRouter)
	r.HandleFunc("/logout", app.authManager.Logout).Methods("GET")
	r.HandleFunc("/user", app.authManager.User).Methods("GET")
	r.HandleFunc("/events", app.eventManager.List).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.GetEvent).Methods("GET")
	//TODO: Проверка на пользователя, отправляющего запрос
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", app.eventManager.CreateEvent).Methods("POST")
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//For test
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello World")
    })

	//Сначала будет вызываться recovery, потом cors, а потом logging
	r.Use(mw.Logging)
	r.Use(mw.CORS)
	//TODO: Убедиться, что достаточно верхней строчки
	/*r.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://bmstusssa.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "DELETE", "POST", "PUT", "OPTIONS"}),
	))
	*/
	r.Use(mw.Recovery)
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
