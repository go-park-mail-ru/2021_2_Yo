package server

import (
	authDelivery "backend/service/auth/delivery/http"
	"backend/service/auth/register"
	authRepository "backend/service/auth/repository/postgres"
	authUseCase "backend/service/auth/usecase"
	register2 "backend/service/user/register"

	userDelivery "backend/service/user/delivery/http"
	userRepository "backend/service/user/repository/postgres"
	userUseCase "backend/service/user/usecase"

	"backend/csrf"
	csrfMiddleware "backend/csrf/middleware"
	csrfRepository "backend/csrf/repository"
	_ "backend/docs"
	eventDelivery "backend/event/delivery/http"
	eventRepository "backend/event/repository/postgres"
	eventUseCase "backend/event/usecase"
	"backend/images"
	imgRepository "backend/images/repository"
	log "backend/logger"
	"backend/middleware"
	"backend/session"
	sessionMiddleware "backend/session/middleware"
	sessionRepository "backend/session/repository"
	"backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const logMessage = "server:"

type App struct {
	AuthManager    *authDelivery.Delivery
	UserManager    *userDelivery.Delivery
	EventManager   *eventDelivery.Delivery
	SessionManager *session.Manager
	CsrfManager    *csrf.Manager
	ImgManager     *images.Manager
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
	redisConnSessions, err := utils.InitRedisDB("redis_db_session")
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	redisConnCSRFTokens, err := utils.InitRedisDB("redis_db_csrf")
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	sessionR := sessionRepository.NewRepository(redisConnSessions)
	sessionM := session.NewManager(*sessionR)

	csrfR := csrfRepository.NewRepository(redisConnCSRFTokens)
	csrfM := csrf.NewManager(*csrfR)

	imgR := imgRepository.NewRepository(db)
	imgM := images.NewManager(*imgR)

	authR := authRepository.NewRepository(db)
	authUC := authUseCase.NewUseCase(authR, []byte(secret))
	authD := authDelivery.NewDelivery(authUC, *sessionM, *csrfM)

	userR := userRepository.NewRepository(db)
	userUC := userUseCase.NewUseCase(userR)
	userD := userDelivery.NewDelivery(userUC, *imgM)

	eventR := eventRepository.NewRepository(db)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		AuthManager:    authD,
		UserManager:    userD,
		EventManager:   eventD,
		SessionManager: sessionM,
		CsrfManager:    csrfM,
		ImgManager:     imgM,
		db:             db,
	}, nil
}

func options(w http.ResponseWriter, r *http.Request) {}

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddleware()
	sessionMW := sessionMiddleware.NewMiddleware(*app.SessionManager)
	csrfMW := csrfMiddleware.NewMiddleware(*app.CsrfManager)

	authRouter := mux.NewRouter()

	CSRFRouter := authRouter.Methods("POST").Subrouter()

	authRouter.HandleFunc("/events/{id:[0-9]+}", app.EventManager.UpdateEvent).Methods("POST")
	authRouter.HandleFunc("/events/{id:[0-9]+}", app.EventManager.DeleteEvent).Methods("DELETE")
	authRouter.HandleFunc("/events", app.EventManager.CreateEvent).Methods("POST")

	CSRFRouter.HandleFunc("/events/{id:[0-9]+}", app.EventManager.UpdateEvent).Methods("POST")
	CSRFRouter.HandleFunc("/events/{id:[0-9]+}", app.EventManager.DeleteEvent).Methods("DELETE")
	CSRFRouter.HandleFunc("/events", app.EventManager.CreateEvent).Methods("POST")
	authRouter.Use(sessionMW.Auth)
	authRouter.Use(mw.GetVars)
	CSRFRouter.Use(csrfMW.CSRF)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(options)

	authSubRouter := r.PathPrefix("/auth").Subrouter()
	register.RegisterHTTPEndpoints(authSubRouter, app, sessionMW)
	userSubRouter := r.PathPrefix("/user").Subrouter()
	register2.RegisterHTTPEndpoints(userSubRouter, app)

	r.Handle("/user", authRouter)

	r.HandleFunc("/user/{id:[0-9]+}", app.UserManager.GetUserById).Methods("GET")

	r.Handle("/user/info", authRouter)
	r.Handle("/user/password", authRouter)

	r.HandleFunc("/events", app.EventManager.GetEventsFromAuthor).Queries("authorid", "{authorid:[0-9]+}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("query", "{query}", "category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("query", "{query}", "category", "{category}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("query", "{query}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("query", "{query}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("category", "{category}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Queries("tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.EventManager.GetEvents).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", app.EventManager.GetEvent).Methods("GET")
	r.Handle("/events/{id:[0-9]+}", authRouter)
	r.Handle("/events", authRouter).Methods("POST")

	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)
	r.Handle("/user/avatar", authRouter)

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
