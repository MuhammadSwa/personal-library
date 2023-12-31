package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN      string
	Port     string
	ENV      string
	LogLevel string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	dbUrl := os.Getenv("DB_DSN")
	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	logLevel := os.Getenv("LOG_LEVEL")

	return &Config{
		DSN:      dbUrl,
		Port:     port,
		ENV:      env,
		LogLevel: logLevel,
	}, nil
}
