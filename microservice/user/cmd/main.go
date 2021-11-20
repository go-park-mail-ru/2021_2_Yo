package main

import (
	"backend/logger"
	proto "backend/microservice/user/proto"
	"backend/microservice/user/repository"
	"backend/utils"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {

	viper.AddConfigPath("../../../configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}

	logLevel := log.DebugLevel
	logger.Init(logLevel)

	db, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(err)
	}

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error(err)
	}

	server := grpc.NewServer()

	userRepositoryService := repository.NewRepository(db)
	proto.RegisterRepositoryServer(server, userRepositoryService)

	log.Info("started user microservice on 8082")
	err = server.Serve(listener)
	if err != nil {
		log.Error("serve troubles")
	}

}
