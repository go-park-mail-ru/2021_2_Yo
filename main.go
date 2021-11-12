package main

import (
	"backend/server"
	"net/smtp"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/joho/godotenv"
)

func env() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func email() {
	from := os.Getenv("EMAIL_ADDR")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := "longhaul2@mail.ru"
	to := []string{toEmail}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := "Golang Email\n"
	body := "Bmstusaaaaaaaa Stepa krutoy!!!!"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("",from,password,host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Check email")
}

func main() {
	log.Info("Main : start")
	env()
	email()
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
