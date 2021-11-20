package main

import (
	protoAuth "backend/microservice/auth/proto"
	sessionRepo "backend/microservice/auth/repository/session"
	userRepo "backend/microservice/auth/repository/user"
	"backend/pkg/logger"
	"backend/pkg/utils"
	log "github.com/sirupsen/logrus"

	"backend/microservice/auth/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

func env() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	env()
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	logLevel := log.DebugLevel
	logger.Init(logLevel)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}

	port := viper.GetString("port")

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
	authListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error(err)
	}

	server := grpc.NewServer()

	authUserRepository := userRepo.NewRepository(postDB)
	authSessionRepository := sessionRepo.NewRepository(redisDB)

	authService := usecase.NewService(authUserRepository, authSessionRepository)
	protoAuth.RegisterAuthServer(server, authService)

	log.Info("started auth microservice on", port)
	err = server.Serve(authListener)
	if err != nil {
		log.Error("serve troubles")
	}

}
