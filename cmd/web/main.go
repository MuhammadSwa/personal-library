package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/muhammadswa/personal-library/config"
	// "github.com/spf13/viper"
	// "github.com/muhammadswa/personal-library/models"
)

// TODO: make people code review this

func main() {
	fmt.Println("Starting the app...")

	fmt.Println("Initializing configuration...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Initializing database...")
	conn, err := InitDatabase(cfg.DSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	fmt.Println("Initializing http server...")
	httpServer := InitHttpServer(conn, cfg.Port)
	err = httpServer.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
