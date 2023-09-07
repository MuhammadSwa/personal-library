package main

import (
	"database/sql"
	"fmt"
)

func InitDatabase(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %w", err)
	}

	return conn, nil
}
