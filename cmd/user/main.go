package main

import (
	"backend/internal/microservice/user/client"
	proto "backend/internal/microservice/user/proto"
	"backend/internal/service/user/repository/postgres"
	"backend/internal/utils"
	log "backend/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
	"github.com/joho/godotenv"
)

func env() {
	// loads values from .env into the system
	if err := godotenv.Load("../../.env"); err != nil {
		log.Error("No .env file found")
	}
}

const logMessage = "cmd:user:"

func main() {
	env()
	logLevel := logrus.DebugLevel
	log.Init(logLevel)
	log.Info(logMessage + "started")

	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}

	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}
	port := viper.GetString("user_port")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}
	server := grpc.NewServer()

	userRepository := postgres.NewRepository(db)
	userClient := client.NewUserService(userRepository)
	proto.RegisterUserServiceServer(server, userClient)

	log.Info("started user microservice on ", port)
	err = server.Serve(listener)
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}
}
