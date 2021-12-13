package app

import (
	protoAuth "backend/internal/microservice/auth/proto"
	"backend/internal/microservice/event/proto"
	userRepository "backend/internal/microservice/user/proto"
	"backend/internal/middleware"
	"backend/internal/notification"
	"backend/internal/register"
	authDelivery "backend/internal/service/auth/delivery/http"
	authUseCase "backend/internal/service/auth/usecase"
	eventDelivery "backend/internal/service/event/delivery/http"
	grpc3 "backend/internal/service/event/repository/grpc"
	eventUseCase "backend/internal/service/event/usecase"
	userDelivery "backend/internal/service/user/delivery/http"
	grpc2 "backend/internal/service/user/repository/grpc"
	userUseCase "backend/internal/service/user/usecase"
	"backend/internal/utils"
	"backend/internal/websocket"
	log "backend/pkg/logger"
	"backend/pkg/prometheus"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	db           *sql.DB
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

	authPort := viper.GetString("auth_port")
	authHost := viper.GetString("auth_host")
	AuthAddr := authHost + ":" + authPort

	grpcConnAuth, err := grpc.Dial(
		AuthAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}

	authClient := protoAuth.NewAuthClient(grpcConnAuth)
	authService := authUseCase.NewUseCase(authClient)
	authD := authDelivery.NewDelivery(authService)

	userPort := viper.GetString("user_port")
	userHost := viper.GetString("user_host")
	userMicroserviceAddr := userHost + ":" + userPort

	userGrpcConn, err := grpc.Dial(userMicroserviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}

	pool := websocket.NewPool()

	SubsNotificator := notification.NewNotificator(pool)

	userRClient := userRepository.NewUserServiceClient(userGrpcConn)
	userR := grpc2.NewRepository(userRClient)
	userUC := userUseCase.NewUseCase(userR)
	userD := userDelivery.NewDelivery(userUC, *SubsNotificator)

	eventPort := viper.GetString("event_port")
	eventHost := viper.GetString("event_host")
	eventMicroserviceAddr := eventHost + ":" + eventPort

	eventGrpcConn, err := grpc.Dial(eventMicroserviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(message+"err = ", err)
		if !opts.Testing {
			return nil, err
		}
	}

	eventRClient := eventGrpc.NewEventServiceClient(eventGrpcConn)
	eventR := grpc3.NewRepository(eventRClient)
	eventUC := eventUseCase.NewUseCase(eventR)
	eventD := eventDelivery.NewDelivery(eventUC)

	return &App{
		Options:      opts,
		AuthManager:  authD,
		UserManager:  userD,
		EventManager: eventD,
		db:           db,
	}, nil
}

func newRouterWithEndpoints(app *App) *mux.Router {
	mw := middleware.NewMiddlewares(app.AuthManager.UseCase)
	mm := prometheus.NewMetricsMiddleware()

	r := mux.NewRouter()
	r.Use(mw.GetVars)
	r.Use(mw.Logging)
	r.Use(mw.CORS)
	r.Use(mw.Recovery)
	r.Use(mm.Metrics)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	authRouter := r.PathPrefix("/auth").Subrouter()
	register.AuthHTTPEndpoints(authRouter, app.AuthManager, mw)
	eventRouter := r.PathPrefix("/events").Subrouter()
	eventRouter.Methods("POST").Subrouter().Use(mw.CSRF)
	register.EventHTTPEndpoints(eventRouter, app.EventManager, mw)
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Methods("POST").Subrouter().Use(mw.CSRF)
	register.UserHTTPEndpoints(userRouter, app.UserManager, app.EventManager, mw)

	r.Handle("/metrics", promhttp.Handler())

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
	_ = r
	res, err1 := app.EventManager.UseCase.GetCities()
	log.Debug(err1)
	log.Debug(res)
	/*
		err := http.ListenAndServe(":"+port, r)
		if err != nil {
			log.Error(message+"err = ", err)
			return err
		}
	*/
	return nil
}
