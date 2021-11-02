package server

import (
	"backend/service/auth"
	authDelivery "backend/service/auth/delivery/http"
	authRepository "backend/service/auth/repository/postgres"
	authUseCase "backend/service/auth/usecase"

	userDelivery "backend/service/user/delivery/http"
	userRepository "backend/service/user/repository/postgres"
	userUseCase "backend/service/user/usecase"

	eventDelivery "backend/service/event/delivery/http"
	eventRepository "backend/service/event/repository/postgres"
	eventUseCase "backend/service/event/usecase"

	"backend/service/csrf"
	csrfRepository "backend/service/csrf/repository"
	"backend/service/image"
	imgRepository "backend/service/image/repository"
	"backend/service/session"
	sessionRepository "backend/service/session/repository"

	_ "backend/docs"
	log "backend/logger"
	"backend/middleware"
	"backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const logMessage = "server:"

type App struct {
	AuthManager    *authDelivery.Delivery
	UserManager    *userDelivery.Delivery
	EventManager   *eventDelivery.Delivery
	SessionManager *session.Manager
	CsrfManager    *csrf.Manager
	ImgManager     *image.Manager
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
	imgM := image.NewManager(*imgR)

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

/*
	mw := middleware.NewMiddleware()
	sessionMW := sessionMiddleware.NewMiddleware(*app.sessionManager)
	csrfMW := csrfMiddleware.NewMiddleware(*app.csrfManager)
	authRouter := mux.NewRouter()
	CSRFRouter := authRouter.Methods("POST").Subrouter()

	authRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	authRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	authRouter.HandleFunc("/events", app.eventManager.CreateEvent).Methods("POST")

	CSRFRouter.HandleFunc("/user/info", app.authManager.UpdateUserInfo).Methods("POST")
	CSRFRouter.HandleFunc("/user/password", app.authManager.UpdateUserPassword).Methods("POST")

	CSRFRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.UpdateEvent).Methods("POST")
	CSRFRouter.HandleFunc("/events/{id:[0-9]+}", app.eventManager.DeleteEvent).Methods("DELETE")
	CSRFRouter.HandleFunc("/events", app.eventManager.CreateEvent).Methods("POST")

	authRouter.Use(sessionMW.Auth)
	authRouter.Use(mw.GetVars)
	CSRFRouter.Use(csrfMW.CSRF)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(options)
	r.Handle("/user", authRouter)
	r.HandleFunc("/user/{id:[0-9]+}", app.authManager.GetUserById).Methods("GET")
	r.Handle("/user/info", authRouter)
	r.Handle("/user/password", authRouter)

	r.HandleFunc("/events", app.eventManager.GetEventsFromAuthor).Queries("authorid", "{authorid:[0-9]+}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("query", "{query}", "category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("query", "{query}", "category", "{category}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("query", "{query}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("query", "{query}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("category", "{category}", "tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("category", "{category}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Queries("tags", "{tags}").Methods("GET")
	r.HandleFunc("/events", app.eventManager.GetEvents).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", app.eventManager.GetEvent).Methods("GET")

	r.Handle("/events/{id:[0-9]+}", authRouter)
	r.Handle("/events", authRouter).Methods("POST")
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)
	r.Handle("/user/avatar", authRouter)

	r.Use(mw.Logging)
	r.Use(mw.CORS)
	r.Use(mw.Recovery)

	return r
*/

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddleware()
	//sessionMW := sessionMiddleware.NewMiddleware(*app.SessionManager)
	//csrfMW := csrfMiddleware.NewMiddleware(*app.CsrfManager)

	r := mux.NewRouter()
	r.Use(mw.Logging)
	r.Use(mw.CORS)
	r.Use(mw.Recovery)
	r.Methods("OPTIONS").HandlerFunc(options)

	authRouter := r.PathPrefix("/auth").Subrouter()
	auth.RegisterHTTPEndpoints(authRouter, app.AuthManager)

	/*
		r.Path("/auth/signup").HandlerFunc(app.AuthManager.SignUp).Methods("POST")
		r.Path("/auth/login").HandlerFunc(app.AuthManager.SignIn).Methods("POST")
		authRouter := r.PathPrefix("/auth").Subrouter()
		authRouter.Use(sessionMW.Auth)
		authRouter.Path("/logout").HandlerFunc(app.AuthManager.Logout).Methods("GET")

		r.Path("/user/{id:[0-9]+}").HandlerFunc(app.UserManager.GetUserById).Methods("GET")
		userRouter := r.PathPrefix("/user").Subrouter()
		userRouter.Use(sessionMW.Auth)
		userRouter.Use(mw.GetVars)
		userRouter.Path("").HandlerFunc(app.UserManager.GetUser).Methods("GET")
		userRouter.Path("info").HandlerFunc(app.UserManager.UpdateUserInfo).Methods("POST")
		userRouter.Path("password").HandlerFunc(app.UserManager.UpdateUserPassword).Methods("POST")
		userRouter.Path("avatar").HandlerFunc(app.UserManager.UpdateUserPhoto).Methods("POST")

	*/

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
