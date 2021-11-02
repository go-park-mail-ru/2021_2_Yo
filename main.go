package main

import (
	"backend/server"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	log.Info("Main : start")

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}

	logLevel := log.DebugLevel
	app, err := server.NewApp(logLevel)
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}
	err = app.Run()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}
}
