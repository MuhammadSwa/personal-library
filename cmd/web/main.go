package main

import (
	_ "github.com/lib/pq"
	"github.com/muhammadswa/personal-library/config"
	"github.com/muhammadswa/personal-library/internal/logger"
)

// TODO: make people code review this

func main() {
	logger.Log.Info().Msg("Starting the app...")

	logger.Log.Info().Msg("Initializing configuration...")
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err, "Error initializing configuration")
	}

	logger.Log.Info().Msg("Initializing database...")
	conn, err := InitDatabase(cfg.DSN)
	if err != nil {
		logger.Fatal(err, "Error initializing database")
	}
	defer conn.Close()

	logger.Log.Info().Msg("Initializing http server...")
	err = InitHttpServer(conn, cfg.Port)
	if err != nil {
		logger.Fatal(err, "Error initializing http server")
	}
}
