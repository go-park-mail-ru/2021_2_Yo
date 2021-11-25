package main

import (
	"backend/server"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func env() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	env()
	log.Info("Main : start")
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main:err = ", err)
		os.Exit(1)
	}
	opts := &server.Options{
		LogLevel: log.DebugLevel,
		Testing:  false,
	}
	app, err := server.NewApp(opts)
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
