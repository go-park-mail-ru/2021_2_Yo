package main

import (
	log "github.com/sirupsen/logrus"
	"backend/logger"
	sessionRepo "backend/microservices/auth/repository/session"
	userRepo "backend/microservices/auth/repository/user"
	protoAuth "backend/microservices/proto/auth"

	"backend/microservices/auth/usecase"
	"backend/utils"
	"net"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func env() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	env()
	viper.AddConfigPath("../../../configs")
	viper.SetConfigName("config")
	logLevel := log.DebugLevel
	logger.Init(logLevel)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}

	//Подключение постгрес
	postDB, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(err)
	}
	//Подключение редис
	redisDB, err := utils.InitRedisDB("redis_db_session")
	if err != nil {
		log.Error(err)
	}

	//Попробую 8081
	authListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Error(err)
	}

	server := grpc.NewServer()

	authUserRepository := userRepo.NewRepository(postDB)
	authSessionRepository := sessionRepo.NewRepository(redisDB)

	authService := usecase.NewService(authUserRepository, authSessionRepository)
	protoAuth.RegisterAuthServer(server,authService)

	log.Info("started auth microservice on 8081")
	err = server.Serve(authListener)
	if err != nil {
		log.Error("serve troubles")
	}

}