package config

import (
	"fmt"
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
	if true {
		dbUrl := "hello"
		fmt.Println(dbUrl)
	}
	port := os.Getenv("PORT")
	return &Config{
		DSN:  dbUrl,
		Port: port,
	}, nil
}
