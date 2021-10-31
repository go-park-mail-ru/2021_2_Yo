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
	"fmt"
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
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	redisConn, err := utils.InitRedisDB()
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	sessionR := sessionRepository.NewRepository(redisConn)
	sessionM := session.NewManager(*sessionR)
	authR := authRepository.NewRepository(db)
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC, *sessionM)
	eventR := eventRepository.NewRepository(db)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		authManager:    authD,
		eventManager:   eventD,
		sessionManager: sessionM,
		db:             db,
	}, nil
}

func options(w http.ResponseWriter, r *http.Request) {}

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddleware()
	sessionMW := sessionMiddleware.NewMiddleware(*app.sessionManager)
	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/logout", app.authManager.Logout).Methods("GET")
	authRouter.HandleFunc("/user", app.authManager.GetUser).Methods("GET")
	authRouter.HandleFunc("/user/info", app.authManager.UpdateUserInfo).Methods("POST")
	authRouter.HandleFunc("/user/password", app.authManager.UpdateUserPassword).Methods("POST")
	authRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	authRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	authRouter.HandleFunc("/events", app.eventManager.CreateEvent).Methods("POST")
	authRouter.Use(mw.GetVars)
	authRouter.Use(sessionMW.Auth)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(options)
	r.HandleFunc("/signup", app.authManager.SignUp).Methods("POST")
	r.HandleFunc("/login", app.authManager.SignIn).Methods("POST")
	r.Handle("/logout", authRouter)
	r.Handle("/user", authRouter)
	r.HandleFunc("/user/{id:[0-9]+}", app.authManager.GetUserById).Methods("GET")
	r.Handle("/user/info", authRouter)
	r.Handle("/user/password", authRouter)
	r.HandleFunc("/events", app.eventManager.List).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.GetEvent).Methods("GET")
	//r.Handle("/events/{id:[0-9]+}", authRouter)
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	r.Handle("/events", authRouter)
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	//Сначала будет вызываться recovery, потом cors, а потом logging
	r.Use(mw.GetVars)
	r.Use(mw.Logging)
	r.Use(mw.CORS)
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
