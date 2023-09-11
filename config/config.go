package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN  string
	Port string
	ENV  string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	dbUrl := os.Getenv("DB_DSN")
	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	return &Config{
		DSN:  dbUrl,
		Port: port,
		ENV:  env,
	}, nil
}
