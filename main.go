package main

import (
	"backend/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Main : start")

	app := server.NewApp()
	app.Run()
}
