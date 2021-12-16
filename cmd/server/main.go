package main

import (
	"backend/internal/app"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"github.com/joho/godotenv"
)

const logMessage = "cmd:server:"

func env() {
	// loads values from .env into the system
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	env()
	log.Info(logMessage + "started")
	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(logMessage+"err = ", err)
		os.Exit(1)
	}
	opts := &app.Options{
		LogLevel: log.DebugLevel,
		Testing:  false,
	}
	application, err := app.NewApp(opts)
	if err != nil {
		log.Error(logMessage+"err = ", err)
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		log.Error(logMessage+"err = ", err)
		os.Exit(1)
	}
}
