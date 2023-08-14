package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func LoanConfig() Config {
	return Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "projectrudi",
			Username: "rudilesmana",
			Password: "rudilesmana2023",
		},
	}
}

func ConnectConfig(config DatabaseConfig) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.Name, config.Username, config.Password,
	)

	connDB, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
		return nil, err
	}

	return connDB, err
}
