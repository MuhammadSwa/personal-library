package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/muhammadswa/personal-lib/internal/config"
)

type application struct {
	cfg *config.Config
}

func main() {
	fmt.Println("Starting the app...")
	app := &application{}

	fmt.Println("Initializing configuration...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}
	app.cfg = cfg

	fmt.Println("Initializing database...")
	conn, err := InitDatabase(app.cfg.DSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	fmt.Println("Initializing http server...")
	err = app.initServer()
	if err != nil {
		log.Fatalln(err)
	}
}
