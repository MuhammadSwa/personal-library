package main

import (
	"database/sql"
	"fmt"
)

func InitDatabase(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %v", err)
	}

	return conn, nil
}
