package server

import (
	authDelivery "backend/service/auth/delivery/http"

	userRepository "backend/microservice/user/proto"
	userDelivery "backend/service/user/delivery/http"
	userUseCase "backend/service/user/usecase"

	eventRepository "backend/microservice/event/proto"
	"backend/register"
	eventDelivery "backend/service/event/delivery/http"
	eventUseCase "backend/service/event/usecase"

	log "backend/logger"
	"backend/middleware"
	"backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	protoAuth "backend/microservice/auth/proto"
	microAuth "backend/service/microservices/auth"
)

const logMessage = "server:"

type App struct {
	AuthManager  *authDelivery.Delivery
	UserManager  *userDelivery.Delivery
	EventManager *eventDelivery.Delivery
	authService  *microAuth.AuthService
	db           *sql.DB
}

func NewApp(logLevel logrus.Level) (*App, error) {
	message := logMessage + "NewApp:"
	log.Init(logLevel)
	log.Info(fmt.Sprintf(message+"started, log level = %s", logLevel))

	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}

	AuthAddr := "localhost:8081"
	grpcConnAuth, err := grpc.Dial(
		AuthAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error("can't connect to grpc")
	}

	authClient := protoAuth.NewAuthClient(grpcConnAuth)
	authService := microAuth.NewService(authClient)
	authD := authDelivery.NewDelivery(authService)

	userMicroserviceAddr := "localhost:8082"
	userGrpcConn, err := grpc.Dial(userMicroserviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err =", err)
	}

	userR := userRepository.NewRepositoryClient(userGrpcConn)
	userUC := userUseCase.NewUseCase(userR)
	userD := userDelivery.NewDelivery(userUC, authService)

	eventMicroserviceAddr := "localhost:8083"
	eventGrpcConn, err := grpc.Dial(eventMicroserviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err =", err)
	}

	eventR := eventRepository.NewRepositoryClient(eventGrpcConn)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		AuthManager:  authD,
		UserManager:  userD,
		EventManager: eventD,
		authService:  &authService,
		db:           db,
	}, nil
}

func options(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	log.Debug("options: ", cookie)
}

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddlewares(*app.authService)

	r := mux.NewRouter()
	r.Use(mw.GetVars)
	r.Use(mw.Logging)
	r.Use(mw.CORS)
	r.Use(mw.Recovery)
	r.Methods("OPTIONS").HandlerFunc(options)

	//TODO: Потом раскоментить и убрать то, что снизу
	//authRouter := r.PathPrefix("/auth").Subrouter()
	//register.AuthHTTPEndpoints(authRouter, app.AuthManager, mw)

	r.HandleFunc("/auth/signup", app.AuthManager.SignUp).Methods("POST")
	r.HandleFunc("/auth/login", app.AuthManager.SignIn).Methods("POST")
	logoutHandlerFunc := http.HandlerFunc(app.AuthManager.Logout)
	r.Handle("/auth/logout", mw.Auth(logoutHandlerFunc))

	eventRouter := r.PathPrefix("/events").Subrouter()
	eventRouter.Methods("POST").Subrouter().Use(mw.CSRF)

	register.EventHTTPEndpoints(eventRouter, app.EventManager, mw)

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Methods("POST").Subrouter().Use(mw.CSRF)
	register.UserHTTPEndpoints(userRouter, app.UserManager, mw)

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
