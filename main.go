package main

import (
	_ "backend/docs"
	"backend/server"

	log "github.com/sirupsen/logrus"
)

//@title BMSTUSA API
//@version 1.0
//@description TP_2021_GO TEAM YO
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

//@host yobmstu.herokuapp.com
//@BasePath /
//@schemes https


func main() {

	log.Info("Main : start")

	app := server.NewApp()
	app.Run()
}

