package main

import (
	_ "backend/docs"
	"backend/server"
	_ "github.com/lib/pq"
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

type user struct {
	id       int
	name     string
	surname  string
	email    string
	password string
}

func main() {
	log.Info("Main : start")
	/*
			connStr := "user=postgres password=password dbname=testDB sslmode=disable"
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Fatal("main : Can't open DB", err)
			}
			defer db.Close()
			log.Println(db.Stats())

		result, err := db.Query(`select * from users`)
		if err != nil{
			panic(err)
		}
		defer result.Close()

		for result.Next(){
			p := user{}
			err := result.Scan(&p.id, &p.name, &p.surname, &p.email, &p.password)
			if err != nil{
				fmt.Println(err)
				continue
			}
			fmt.Println(p.id, p.name, p.surname, p.email, p.password)
		}
	*/

	/*
	app := server.NewApp()
	app.Run()

	 */
	server.NewApp()
}
