package main

import (
	_ "github.com/lib/pq"
	"github.com/muhammadswa/personal-library/config"

	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO: make people code review this

func main() {

	log.Info().Msg("Starting the app...")

	log.Info().Msg("Initializing configuration...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing configuration")
		// log.Fatalln(err)
	}

	log.Info().Msg("Initializing database...")
	conn, err := InitDatabase(cfg.DSN)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing configuration")
	}
	defer conn.Close()

	log.Info().Msg("Initializing http server...")
	httpServer := InitHttpServer(conn, cfg.Port)
	err = httpServer.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing configuration")
	}
}
