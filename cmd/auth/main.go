package main

import (
	protoAuth "backend/internal/microservice/auth/proto"
	sessionRepo "backend/internal/microservice/auth/repository/session"
	userRepo "backend/internal/microservice/auth/repository/user"
	"backend/internal/microservice/auth/usecase"
	"backend/internal/utils"
	log "backend/pkg/logger"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

const logMessage = "cmd:auth:"

func main() {
	logLevel := logrus.DebugLevel
	log.Init(logLevel)
	log.Info(logMessage + "started")

	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(logMessage+"err = ", err)
		os.Exit(1)
	}

	port := viper.GetString("auth_port")

	postDB, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(logMessage+"err = ", err)
	}
	redisDB, err := utils.InitRedisDB()
	if err != nil {
		log.Error(logMessage+"err = ", err)
	}
	authListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error(logMessage+"err = ", err)
	}

	server := grpc.NewServer()

	authUserRepository := userRepo.NewRepository(postDB)
	authSessionRepository := sessionRepo.NewRepository(redisDB)

	authService := usecase.NewService(authUserRepository, authSessionRepository)
	protoAuth.RegisterAuthServer(server, authService)

	log.Info("started auth microservice on ", port)
	err = server.Serve(authListener)
	if err != nil {
		log.Error("serve troubles")
	}

}
