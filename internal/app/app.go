package app

import (
	protoAuth "backend/internal/microservice/auth/proto"
	eventRepository "backend/internal/microservice/event/proto"
	userRepository "backend/internal/microservice/user/proto"
	"backend/internal/middleware"
	"backend/internal/register"
	authDelivery "backend/internal/service/auth/delivery/http"
	authUseCase "backend/internal/service/auth/usecase"
	eventDelivery "backend/internal/service/event/delivery/http"
	eventGrpc "backend/internal/service/event/repository/grpc"
	eventUseCase "backend/internal/service/event/usecase"
	"backend/internal/service/notification/delivery/websocket"
	"backend/internal/service/notification/repository/postgres"
	userDelivery "backend/internal/service/user/delivery/http"
	userGrpc "backend/internal/service/user/repository/grpc"
	userUseCase "backend/internal/service/user/usecase"
	"backend/internal/utils"
	log "backend/pkg/logger"
	"backend/pkg/notificator"
	//"backend/pkg/prometheus"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const logMessage = "server:"

type Options struct {
	LogLevel logrus.Level
	Testing  bool
}

type App struct {
	Options      *Options
	AuthManager  *authDelivery.Delivery
	UserManager  *userDelivery.Delivery
	EventManager *eventDelivery.Delivery
	wsPool       *websocket.Pool
	db           *sql.DB
}

func getGrpcAddress(portKey string, hostKey string) string {
	port := viper.GetString(portKey)
	host := viper.GetString(hostKey)
	return host + ":" + port
}

func NewApp(opts *Options) (*App, error) {
	message := logMessage + "NewApp:"
	log.Init(opts.LogLevel)
	log.Info(fmt.Sprintf(message+"started, log level = %s", opts.LogLevel))

	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}

	grpcConnAuth, err := grpc.Dial(getGrpcAddress("auth_port", "auth_host"), grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}
	userGrpcConn, err := grpc.Dial(getGrpcAddress("user_port", "user_host"), grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}
	eventGrpcConn, err := grpc.Dial(getGrpcAddress("event_port", "event_host"), grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}

	authClient := protoAuth.NewAuthClient(grpcConnAuth)
	userRClient := userRepository.NewUserServiceClient(userGrpcConn)
	eventRClient := eventRepository.NewEventServiceClient(eventGrpcConn)

	notificationR := postgres.NewRepository(db)
	userR := userGrpc.NewRepository(userRClient)
	eventR := eventGrpc.NewRepository(eventRClient)

	authService := authUseCase.NewUseCase(authClient)
	userUC := userUseCase.NewUseCase(userR)
	eventUC := eventUseCase.NewUseCase(eventR)

	pool := websocket.NewPool()
	notificationManager := notificator.NewNotificator(pool, notificationR, userR, eventR)

	authD := authDelivery.NewDelivery(authService)
	userD := userDelivery.NewDelivery(userUC, notificationManager)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		Options:      opts,
		AuthManager:  authD,
		UserManager:  userD,
		EventManager: eventD,
		wsPool:       pool,
		db:           db,
	}, nil
}

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddlewares(app.AuthManager.UseCase)
	//mm := prometheus.NewMetricsMiddleware()

	r := mux.NewRouter()
	rApi := r.PathPrefix("/api").Subrouter()
	rApi.Use(mw.GetVars)
	rApi.Use(mw.Logging)
	rApi.Use(mw.CORS)
	rApi.Use(mw.Recovery)
	//r.Use(mm.Metrics)
	rApi.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	authRouter := rApi.PathPrefix("/auth").Subrouter()
	register.AuthHTTPEndpoints(authRouter, app.AuthManager, mw)
	eventRouter := rApi.PathPrefix("/events").Subrouter()
	eventRouter.Methods("POST").Subrouter().Use(mw.CSRF)
	register.EventHTTPEndpoints(eventRouter, app.EventManager, mw)
	userRouter := rApi.PathPrefix("/user").Subrouter()
	userRouter.Methods("POST").Subrouter().Use(mw.CSRF)
	register.UserHTTPEndpoints(userRouter, app.UserManager, app.EventManager, mw)
	r.HandleFunc("/ws", app.wsPool.WebsocketHandler).Methods("GET")

	//r.Handle("/metrics", promhttp.Handler())

	return r
}

func (app *App) Run() error {
	if app.db != nil {
		defer app.db.Close()
	}
	message := logMessage + "Run:"
	log.Info(message + "start")
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("bmstusa_port")
	}
	log.Info(message+"port = ", port)
	if app.Options.Testing {
		port = "test port"
	}
	r := newRouterWithEndpoints(app)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error(message+"err = ", err)
		return err
	}
	return nil
}
