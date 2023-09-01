package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN  string
	Port string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	dbUrl := os.Getenv("DB_DSN")
	port := os.Getenv("PORT")
	return &Config{
		DSN:  dbUrl,
		Port: port,
	}, nil
}
