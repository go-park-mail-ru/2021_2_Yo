package main

import (
	_ "backend/docs"
	"backend/server"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
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
	app, err := server.NewApp()
	if err != nil {
		log.Error("Main : NewApp error = ", err)
		os.Exit(1)
	}
	err = app.Run()
	if err != nil {
		log.Error("Main : Run error = ", err)
		os.Exit(1)
	}
}
