package main

import (
	proto "backend/microservice/event/proto"
	repository "backend/microservice/event/repository"
	log "backend/pkg/logger"
	"backend/utils"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

const logMessage = "microservice:event:"

func main() {

	logLevel := logrus.DebugLevel
	log.Init(logLevel)

	log.Info(logMessage + "started")

	viper.AddConfigPath("../../../configs")
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

	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}

	server := grpc.NewServer()

	eventRepositoryService := repository.NewRepository(db)
	proto.RegisterRepositoryServer(server, eventRepositoryService)

	log.Info("started event microservice on 8083")
	err = server.Serve(listener)
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}

}
