package main

import (
	"backend/microservice/event/client"
	proto "backend/microservice/event/proto"
	log "backend/pkg/logger"
	"backend/pkg/utils"
	repository "backend/service/event/repository/postgres"
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

	viper.AddConfigPath("../config")
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

	port := viper.GetString("event_port")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}

	server := grpc.NewServer()

	eventRepository := repository.NewRepository(db)
	eventService := client.NewEventService(eventRepository)
	proto.RegisterEventServiceServer(server, eventService)

	log.Info("started event microservice on ", port)
	err = server.Serve(listener)
	if err != nil {
		log.Error(logMessage+"err =", err)
		os.Exit(1)
	}

}
